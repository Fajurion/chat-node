package pipe

import (
	"chat-node/util"
	"log"

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

	log.Println(nodeList)

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
