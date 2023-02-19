package pipe

import (
	"chat-node/util"
	"log"
	"os"

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

var Nodes map[int64]Node = make(map[int64]Node)

func queryNodes() (bool, bool) {

	res, err := util.PostRequest("/node/list", fiber.Map{
		"token": util.NODE_TOKEN,
	})

	if err != nil {
		return true, false
	}

	if !res["success"].(bool) {
		return true, false
	}

	if res["nodes"] == nil {
		return false, true
	}

	nodeList := res["nodes"].([]interface{})

	for _, node := range nodeList {

		n := node.(map[string]interface{})
		Nodes[int64(n["id"].(float64))] = Node{
			ID:     int64(n["id"].(float64)),
			Token:  n["token"].(string),
			App:    uint(n["app"].(float64)),
			Domain: n["domain"].(string),
		}
	}

	return false, false
}

var CurrentNode Node

func queryNode() {

	res, err := util.PostRequest("/node/this", fiber.Map{
		"id":    util.NODE_ID,
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
