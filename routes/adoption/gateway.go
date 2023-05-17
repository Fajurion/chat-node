package adoption

import (
	"chat-node/util/requests"

	integration "fajurion.com/node-integration"
	"github.com/Fajurion/pipes"
	"github.com/Fajurion/pipes/connection"
	"github.com/Fajurion/pipes/receive"
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

			// Adopt node
			node, err := receive.ReceiveWSAdoption(token)
			if err != nil {
				return requests.InvalidRequest(c)
			}

			// Set the token as a local variable
			c.Locals("ws", true)
			c.Locals("node", node)
			return c.Next()
		}

		return c.SendStatus(fiber.StatusUpgradeRequired)
	})

	router.Get("/", websocket.New(ws))

}

func ws(conn *websocket.Conn) {
	node := conn.Locals("node").(pipes.Node)

	defer func() {

		// Disconnect node
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
