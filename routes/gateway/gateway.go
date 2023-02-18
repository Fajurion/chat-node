package gateway

import (
	"chat-node/bridge"
	"chat-node/handler"
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

			if len(token) == 0 {
				return c.SendStatus(fiber.StatusUnauthorized)
			}

			// Check if the token is valid
			tk := bridge.CheckToken(token)
			if tk.UserID == 0 {
				return c.SendStatus(fiber.StatusBadRequest)
			}

			bridge.RemoveToken(token)

			// Set the token as a local variable
			c.Locals("ws", true)
			c.Locals("token", token)
			c.Locals("tk", tk)
			return c.Next()
		}

		return c.SendStatus(fiber.StatusUpgradeRequired)
	})

	router.Get("/", websocket.New(ws))
}

type Message struct {
	Action string                 `json:"action"`
	Data   map[string]interface{} `json:"data"`
}

func ws(conn *websocket.Conn) {
	tk := conn.Locals("tk").(bridge.ConnectionToken)
	token := conn.Locals("token").(string)

	bridge.AddClient(conn, tk.UserID, token, tk.Session)
	defer bridge.Remove(tk.UserID, token)

	for {
		// Read message as text
		mtype, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}

		log.Println(mtype)

		// Broadcast msg
		if mtype == websocket.TextMessage {

			// Unmarshal the event
			var message Message
			err := sonic.UnmarshalString(string(msg), &message)
			if err != nil {
				return
			}

			// Handle the event
			if !handler.Handle(handler.Message{
				Client: bridge.Get(tk.UserID, token),
				Data:   message.Data,
				Action: message.Action,
			}) {
				return
			}

		} else {
			bridge.Remove(tk.UserID, token)
		}
	}
}
