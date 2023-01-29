package bridge

import (
	"chat-node/util"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

type Client struct {
	Conn    *websocket.Conn
	ID      int64
	Session string
	End     time.Time
}

func (c *Client) IsExpired() bool {
	return c.End.Before(time.Now())
}

// User ID -> Token -> Client
var Connections map[int64]map[string]Client = make(map[int64]map[string]Client)

func AddClient(conn *websocket.Conn, id int64, token string, session string) {
	log.Println("New connection", token)

	if Connections[id] == nil {
		Connections[id] = make(map[string]Client)
	}

	Connections[id][token] = Client{
		Conn:    conn,
		ID:      id,
		Session: session,
	}
}

func Remove(id int64, token string) {
	log.Println("Connection closed", token)
	Connections[id][token].Conn.Close()

	// Send to server
	util.PostRequest("/node/disconnect", fiber.Map{
		"node_token": util.NODE_TOKEN,
		"token":      Connections[id][token].Session,
	})

	delete(Connections[id], token)

	if len(Connections[id]) == 0 {
		delete(Connections, id)
	}
}

func Broadcast(msg []byte) {
	for _, clients := range Connections {
		for _, client := range clients {
			SendMessage(client.Conn, msg)
		}
	}
}

func Send(id int64, msg []byte) {
	for _, client := range Connections[id] {
		SendMessage(client.Conn, msg)
	}
}

func SendMessage(conn *websocket.Conn, msg []byte) {
	conn.WriteMessage(websocket.TextMessage, msg)
}