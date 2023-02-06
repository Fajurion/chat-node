package pipe

import (
	"chat-node/util"

	"github.com/gofiber/fiber/v2"
)

type Node struct {
	ID     int64  `json:"id"`
	App    uint   `json:"app"`
	Token  string `json:"token"`
	Domain string `json:"domain"`
}

func (n Node) GetWebSocket() string {
	return "ws://" + n.Domain + "/adoption"
}

var Nodes map[int64]Node

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

	nodeList := res["nodes"].([]Node)

	for _, node := range nodeList {
		Nodes[node.ID] = node
	}

	return false, false
}
