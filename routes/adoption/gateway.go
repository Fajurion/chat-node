package adoption

import (
	"chat-node/util"
	"log"

	integration "fajurion.com/node-integration"
	"github.com/Fajurion/pipes"
	"github.com/Fajurion/pipes/connection"
	"github.com/Fajurion/pipes/receive"
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func SetupRoutes(router fiber.Router) {

	router.Post("/socketless", socketless)

	// Inject a middleware to check if the request is a websocket upgrade request
	router.Use("/", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {

			// Check if the request has a token
			token := c.Get("Sec-WebSocket-Protocol")

			// Parse request
			var req connection.AdoptionRequest
			if err := sonic.Unmarshal([]byte(token), &req); err != nil {
				return c.SendStatus(fiber.StatusBadRequest)
			}

			// Check if the token is valid
			if util.NODE_TOKEN != req.Token {
				return c.SendStatus(fiber.StatusUnauthorized)
			}

			pipes.AddNode(req.Adopting)

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
	node := conn.Locals("node").(pipes.Node)

	// Check if event connection is already open
	if !connection.ExistsWS(node.ID) {
		log.Println("Building outgoing event stream to node", node.ID)
		go connection.ConnectWS(node)
	}

	log.Printf("Incoming event stream of node %s connected. \n", node.ID)
	defer func() {

		// Close connection
		log.Printf("Incoming event stream of node %s disconnected. \n", node.ID)

		connection.RemoveWS(node.ID)
		integration.ReportOffline(node)
		conn.Close()
	}()

	for {
		// Read message as text
		mtype, msg, err := conn.ReadMessage()
		if err != nil {
			return
		}

		if mtype == websocket.TextMessage {

			// Pass message to pipes
			receive.ReceiveWS(msg)
		}
	}

}
