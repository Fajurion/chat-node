package adoption

import (
	"chat-node/pipe"
	"chat-node/util"
	"log"
	"strconv"
	"strings"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func SetupRoutes(router fiber.Router) {

	router.Post("/initialize", initialize)

	// Inject a middleware to check if the request is a websocket upgrade request
	router.Use("/", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {

			// Check if the request has a token
			token := c.Get("Sec-WebSocket-Protocol")

			if len(token) == 0 {
				return c.SendStatus(fiber.StatusUnauthorized)
			}

			// Check if the token is valid
			args := strings.Split(token, "_")

			if len(args) != 3 {
				return c.SendStatus(fiber.StatusUnauthorized)
			}

			id, err := strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				return c.SendStatus(fiber.StatusUnauthorized)
			}

			if args[0] != util.NODE_TOKEN || args[2] != pipe.Nodes[id].Token {
				return c.SendStatus(fiber.StatusUnauthorized)
			}

			// Set the token as a local variable
			c.Locals("ws", true)
			c.Locals("node", pipe.Nodes[id])
			return c.Next()
		}

		return c.SendStatus(fiber.StatusUpgradeRequired)
	})

	router.Get("/", websocket.New(ws))

}

func ws(conn *websocket.Conn) {
	node := conn.Locals("node").(pipe.Node)

	// Check if event connection is already open
	if pipe.NodeConnections[node.ID] == nil {
		go pipe.ConnectToNode(node)
	}

	log.Printf("Incoming event stream of node %d connected. \n", node.ID)
	defer func() {

		// Close connection
		log.Printf("Incoming event stream of node %d disconnected. \n", node.ID)

		pipe.ReportOffline(node)
		pipe.NodeConnections[node.ID].Close(websocket.CloseInternalServerErr, "incoming.closed")
		delete(pipe.NodeConnections, node.ID)

		conn.Close()
	}()

	for {
		// Read message as text
		mtype, msg, err := conn.ReadMessage()
		if err != nil {
			return
		}

		if mtype == websocket.TextMessage {

			// Parse message
			var message pipe.Message
			if err := sonic.Unmarshal(msg, &message); err != nil {
				return
			}

			// Send message to node
			log.Println("Received message:", message)

		}
	}

}
