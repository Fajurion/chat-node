package adoption

import (
	"chat-node/pipe"
	"chat-node/util"
	"log"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func SetupRoutes(router fiber.Router) {

	// Inject a middleware to check if the request is a websocket upgrade request
	router.Use("/", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {

			// Check if the request has a token
			token := c.Get("Sec-WebSocket-Protocol")

			// Parse request
			var req pipe.AdoptionRequest
			if err := sonic.Unmarshal([]byte(token), &req); err != nil {
				return c.SendStatus(fiber.StatusBadRequest)
			}

			// Check if the token is valid
			if util.NODE_TOKEN != req.Token {
				return c.SendStatus(fiber.StatusUnauthorized)
			}

			pipe.AddNode(req.Adopting)

			// Set the token as a local variable
			c.Locals("ws", true)
			c.Locals("node", req.Adopting)
			return c.Next()
		}

		return c.SendStatus(fiber.StatusUpgradeRequired)
	})

	router.Get("/", websocket.New(ws))

}

func ws(conn *websocket.Conn) {
	node := conn.Locals("node").(pipe.Node)

	// Check if event connection is already open
	if !pipe.ConnectionExists(node.ID) {
		go pipe.ConnectToNode(node)
	}

	log.Printf("Incoming event stream of node %d connected. \n", node.ID)
	defer func() {

		// Close connection
		log.Printf("Incoming event stream of node %d disconnected. \n", node.ID)

		pipe.Offline(node)
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
