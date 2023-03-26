package pipe

import (
	"chat-node/util"
	"log"
	"os"
	"os/exec"

	"github.com/gofiber/fiber/v2"
)

func Create() {

	log.Println("Creating pipe...")

	// Get current node
	log.Printf("Getting current node info... (ID: %d)", util.NODE_ID)
	queryNode()

	errTitle := exec.Command("cmd", "/C", "title", CurrentNode.Domain).Run()
	if errTitle != nil {
		log.Println("Failed to set window title.")
	}

	log.Println("Current node info:", CurrentNode)

	// Tell server new status
	res := updateStatus()

	// Get all nodes
	log.Println("Connecting to other nodes..")
	err, solo := parseNodes(res)

	if err {
		log.Println("Backend is currently offline!")
		os.Exit(1)
	}

	if solo {
		log.Println("This node is currently running solo mode.")
	} else {

		// Connect to all nodes
		IterateNodes(func(_ int64, node Node) bool {
			log.Println("Connecting to node " + node.Domain + "...")
			ConnectToNode(node)
			return true
		})
	}

	log.Println("Updated node status to online.")
}

// updateStatus updates the node status to online
func updateStatus() map[string]interface{} {
	res, err := util.PostRequest("/node/status/online", fiber.Map{
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

	return res
}
