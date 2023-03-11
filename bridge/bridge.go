package bridge

import (
	"chat-node/pipe"
	"time"

	"github.com/bytedance/sonic"
	"github.com/cornelk/hashmap"
	"github.com/gofiber/websocket/v2"
)

type Client struct {
	Conn     *websocket.Conn
	ID       int64
	Session  uint64
	Username string
	Tag      string
	End      time.Time
}

func (c *Client) SendEvent(event pipe.Event) {

	event.Sender = c.ID
	msg, err := sonic.Marshal(event)
	if err != nil {
		return
	}

	SendMessage(c.Conn, msg)
}

func (c *Client) IsExpired() bool {
	return c.End.Before(time.Now())
}

// User ID -> Token -> Client
var Connections = hashmap.New[int64, *hashmap.Map[uint64, Client]]()

func AddClient(conn *websocket.Conn, id int64, session uint64, username string, tag string) {

	if _, ok := Connections.Get(id); !ok {
		Connections.Insert(id, hashmap.New[uint64, Client]())
	}

	clients, _ := Connections.Get(id)

	clients.Insert(session, Client{
		Conn:     conn,
		ID:       id,
		Session:  session,
		Username: username,
		Tag:      tag,
	})

	Connections.Set(id, clients)
}

func Remove(id int64, session uint64) {

	clients, ok := Connections.Get(id)

	if !ok {
		return
	}

	clients.Del(session)

	if clients.Len() == 0 {
		Connections.Del(id)
	} else {
		Connections.Set(id, clients)
	}
}

func Send(id int64, msg []byte) {
	clients, ok := Connections.Get(id)

	if !ok {
		return
	}

	clients.Range(func(session uint64, client Client) bool {

		SendMessage(client.Conn, msg)
		return true
	})
}

func SendSession(id int64, session uint64, msg []byte) {
	clients, _ := Connections.Get(id)
	client, _ := clients.Get(session)

	SendMessage(client.Conn, msg)
}

func SendMessage(conn *websocket.Conn, msg []byte) {
	conn.WriteMessage(websocket.TextMessage, msg)
}

func ExistsConnection(id int64, session uint64) bool {
	clients, ok := Connections.Get(id)
	if !ok {
		return false
	}

	_, ok = clients.Get(session)
	return ok
}

func Get(id int64, session uint64) *Client {
	clients, _ := Connections.Get(id)
	client, _ := clients.Get(session)

	return &client
}

func GetConnections(id int64) int {
	clients, ok := Connections.Get(id)
	if !ok {
		return 0
	}

	return clients.Len()
}
