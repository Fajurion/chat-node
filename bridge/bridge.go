package bridge

import (
	"chat-node/util"
	"log"
	"time"

	"github.com/cornelk/hashmap"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

type Client struct {
	Conn     *websocket.Conn
	ID       int64
	Session  uint
	Username string
	Tag      string
	End      time.Time
}

func (c *Client) IsExpired() bool {
	return c.End.Before(time.Now())
}

// User ID -> Token -> Client
var Connections = hashmap.New[int64, *hashmap.Map[string, Client]]()

func AddClient(conn *websocket.Conn, id int64, token string, session uint) {
	log.Println("New connection", token)

	if _, ok := Connections.Get(id); !ok {
		Connections.Insert(id, hashmap.New[string, Client]())
	}

	clients, _ := Connections.Get(id)

	clients.Insert(token, Client{
		Conn:    conn,
		ID:      id,
		Session: session,
	})

	Connections.Set(id, clients)
}

func Remove(id int64, token string) {
	log.Println("Connection closed", token)
	clients, _ := Connections.Get(id)
	client, _ := clients.Get(token)

	// Send to server
	util.PostRequest("/node/disconnect", fiber.Map{
		"node_token": util.NODE_TOKEN,
		"token":      client.Session,
	})

	clients.Del(token)

	if clients.Len() == 0 {
		Connections.Del(id)
	} else {
		Connections.Set(id, clients)
	}
}

func Send(id int64, msg []byte) {
	clients, _ := Connections.Get(id)

	clients.Range(func(key string, client Client) bool {

		SendMessage(client.Conn, msg)
		return true
	})
}

func SendMessage(conn *websocket.Conn, msg []byte) {
	conn.WriteMessage(websocket.TextMessage, msg)
}

func Get(id int64, token string) Client {
	clients, _ := Connections.Get(id)
	client, _ := clients.Get(token)

	return client
}
