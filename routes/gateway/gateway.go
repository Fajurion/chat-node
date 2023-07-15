package gateway

import (
	"chat-node/bridge"
	"chat-node/database/fetching"
	"chat-node/handler"
	"chat-node/handler/account"
	"chat-node/service"
	"log"

	"github.com/Fajurion/pipes/adapter"
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
			if tk.UserID == "" {
				return c.SendStatus(fiber.StatusBadRequest)
			}

			if bridge.ExistsConnection(tk.UserID, tk.Session) {
				return c.SendStatus(fiber.StatusConflict)
			}

			bridge.RemoveToken(token)

			// Set the token as a local variable
			c.Locals("ws", true)
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

	bridge.AddClient(conn, tk.UserID, tk.Session, tk.Username, tk.Tag)
	defer func() {

		// Update status
		account.UpdateStatus(bridge.Get(tk.UserID, tk.Session), fetching.StatusOffline, "", false)

		// Remove the connection from the bridge
		bridge.Remove(tk.UserID, tk.Session)
	}()

	if !service.User(bridge.Get(tk.UserID, tk.Session)) {
		log.Println("Something's wrong with the user")
		return
	}

	// Add adapter for pipes
	adapter.AdaptWS(adapter.Adapter{
		ID: tk.UserID,
		Receive: func(c *adapter.Context) error {
			return conn.WriteMessage(websocket.TextMessage, c.Message)
		},
	})

	for {
		// Read message as text
		_, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}

		// Unmarshal the event
		var message Message
		err = sonic.UnmarshalString(string(msg), &message)
		if err != nil {
			return
		}

		// Handle the event
		if !handler.Handle(handler.Message{
			Client: bridge.Get(tk.UserID, tk.Session),
			Data:   message.Data,
			Action: message.Action,
		}) {
			return
		}
	}
}
