package pipe

import (
	"chat-node/util"
	"log"
	"os"

	"github.com/cornelk/hashmap"
	"github.com/gofiber/fiber/v2"
)

type Node struct {
	ID     int64  `json:"id"`
	Token  string `json:"token"`
	App    uint   `json:"app"`
	Domain string `json:"domain"`
}

func (n Node) GetWebSocket() string {
	return "ws://" + n.Domain + "/adoption"
}

var nodes = hashmap.New[int64, Node]()

func parseNodes(res map[string]interface{}) (bool, bool) {

	if res["nodes"] == nil {
		return false, true
	}

	nodeList := res["nodes"].([]interface{})

	for _, node := range nodeList {

		n := node.(map[string]interface{})
		nodes.Insert(int64(n["id"].(float64)), Node{
			ID:     int64(n["id"].(float64)),
			Token:  n["token"].(string),
			App:    uint(n["app"].(float64)),
			Domain: n["domain"].(string),
		})
	}

	return false, false
}

var CurrentNode Node

func queryNode() {

	res, err := util.PostRequest("/node/this", fiber.Map{
		"node":  util.NODE_ID,
		"token": util.NODE_TOKEN,
	})

	if err != nil {
		log.Println("Backend is currently offline!")
		os.Exit(1)
	}

	if !res["success"].(bool) {
		log.Println("This node may not be registered..")
		os.Exit(1)
	}

	n := res["node"].(map[string]interface{})

	CurrentNode = Node{
		ID:     int64(n["id"].(float64)),
		Token:  n["token"].(string),
		App:    uint(n["app"].(float64)),
		Domain: n["domain"].(string),
	}
}

func GetNode(id int64) *Node {

	// Get node
	node, ok := nodes.Get(id)
	if !ok {
		return nil
	}

	return &node
}

func AddNode(node Node) {
	nodes.Insert(node.ID, node)
}

// IterateConnections iterates over all connections. If the callback returns false, the iteration stops.
func IterateNodes(callback func(int64, Node) bool) {
	nodes.Range(callback)
}
