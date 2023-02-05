package gateway

import (
	"chat-node/bridge"
	"chat-node/bridge/conversation"
	"chat-node/pipe"
	"chat-node/pipe/send"

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

		// Broadcast msg
		if mtype == websocket.TextMessage {

			// Unmarshal the event
			var message pipe.Message
			err := sonic.UnmarshalString(string(msg), &message)
			if err != nil {
				bridge.Remove(tk.UserID, token)
				continue
			}

			// Check if the event is valid
			event := message.Event
			channel := message.Channel

			if channel.Sender != tk.UserID || len(event.Name) == 0 || len(event.Data) == 0 {
				bridge.Remove(tk.UserID, token)
				continue
			}

			if !channel.IsValid(event) {
				bridge.Remove(tk.UserID, token)
				continue
			}

			// Check if the user is in the project
			if channel.IsProject() {
				if conversation.GetProject(channel.Target[0]).Members[tk.UserID] == 0 {
					bridge.Remove(tk.UserID, token)
					continue
				}
			}

			// Send the event to the pipe
			send.Pipe(channel, msg, event)

		} else {
			bridge.Remove(tk.UserID, token)
		}
	}
}
