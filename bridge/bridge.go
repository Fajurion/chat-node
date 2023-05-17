package bridge

import (
	"chat-node/database"
	"chat-node/database/fetching"
	"chat-node/util"
	"log"
	"time"

	"github.com/Fajurion/pipes"
	"github.com/bytedance/sonic"
	"github.com/cornelk/hashmap"
	"github.com/gofiber/websocket/v2"
)

type Client struct {
	Conn     *websocket.Conn
	ID       string
	Session  string
	Room     string // Livekit room
	Username string
	Tag      string
	End      time.Time
}

func (c *Client) SendEvent(event pipes.Event) {

	log.Println(event.Name)

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
var Connections = hashmap.New[string, *hashmap.Map[string, Client]]()

func AddClient(conn *websocket.Conn, id string, session string, username string, tag string) {

	if _, ok := Connections.Get(id); !ok {
		Connections.Insert(id, hashmap.New[string, Client]())
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

func Remove(id string, session string) {

	database.DBConn.Model(&fetching.Session{}).Where("id = ?", session).Update("last_fetch", time.Now().UnixMilli())
	util.PostRequest("/node/disconnect", map[string]interface{}{
		"node":    util.NODE_ID,
		"token":   util.NODE_TOKEN,
		"session": session,
	})

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

func Send(id string, msg []byte) {
	clients, ok := Connections.Get(id)

	if !ok {
		return
	}

	clients.Range(func(session string, client Client) bool {

		SendMessage(client.Conn, msg)
		return true
	})
}

func SendSession(id string, session string, msg []byte) {
	clients, _ := Connections.Get(id)
	client, _ := clients.Get(session)

	SendMessage(client.Conn, msg)
}

func SendMessage(conn *websocket.Conn, msg []byte) {
	conn.WriteMessage(websocket.TextMessage, msg)
}

func ExistsConnection(id string, session string) bool {
	clients, ok := Connections.Get(id)
	if !ok {
		return false
	}

	_, ok = clients.Get(session)
	return ok
}

func Get(id string, session string) *Client {
	clients, _ := Connections.Get(id)
	client, _ := clients.Get(session)

	return &client
}

func GetConnections(id string) int {
	clients, ok := Connections.Get(id)
	if !ok {
		return 0
	}

	return clients.Len()
}
