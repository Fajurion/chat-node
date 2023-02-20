package pipe

import (
	"chat-node/util"
	"context"
	"log"

	"github.com/bytedance/sonic"
	"github.com/cornelk/hashmap"
	"github.com/gofiber/fiber/v2"
	"nhooyr.io/websocket"
)

var nodeConnections = hashmap.New[int64, *websocket.Conn]()

type AdoptionRequest struct {
	Token    string `json:"tk"`
	Adopting Node   `json:"adpt"`
}

func ConnectToNode(node Node) {

	// Marshal current node
	nodeBytes, err := sonic.Marshal(AdoptionRequest{
		Token:    node.Token,
		Adopting: CurrentNode,
	})
	if err != nil {
		return
	}

	// Connect to node
	c, _, err := websocket.Dial(context.Background(), node.GetWebSocket(), &websocket.DialOptions{
		Subprotocols: []string{string(nodeBytes)},
	})

	if err != nil {
		return
	}

	// Add connection to map
	nodeConnections.Insert(node.ID, c)

	log.Printf("Outgoing event stream to node %d connected.", node.ID)
}

func Offline(node Node) {

	// Check if connection exists
	if !ConnectionExists(node.ID) {
		return
	}

	_, err := util.PostRequest("/node/status/offline", fiber.Map{
		"token": node.Token,
	})

	if err != nil {
		log.Println("Failed to report offline status. Is the backend online?")
	}

	connection := GetConnection(node.ID)
	connection.Close(websocket.StatusNormalClosure, "node.offline")

	nodeConnections.Del(node.ID)
}

func ConnectionExists(node int64) bool {

	// Check if connection exists
	_, ok := nodeConnections.Get(node)
	if !ok {
		return false
	}

	return true
}

func GetConnection(node int64) *websocket.Conn {

	// Check if connection exists
	connection, ok := nodeConnections.Get(node)
	if !ok {
		return nil
	}

	return connection
}

// Range calls f sequentially for each key and value present in the map. If f returns false, range stops the iteration.
func IterateConnections(f func(key int64, value *websocket.Conn) bool) {
	nodeConnections.Range(f)
}
