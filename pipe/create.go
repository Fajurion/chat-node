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
	log.Println("Getting current node info...")
	queryNode()

	errTitle := exec.Command("cmd", "/C", "title", CurrentNode.Domain).Run()
	if errTitle != nil {
		log.Println("Failed to set window title.")
	}

	log.Println("Current node info:", CurrentNode)

	// Tell server new status
	updateStatus()

	// Get all nodes
	log.Println("Connecting to other nodes..")
	err, solo := queryNodes()

	if err {
		log.Println("Backend is currently offline!")
		os.Exit(1)
	}

	if solo {
		log.Println("This node is currently running solo mode.")
	} else {

		// Connect to all nodes
		for _, node := range Nodes {
			log.Println("Connecting to node " + node.Domain + "...")
			ConnectToNode(node)
		}
	}

	log.Println("Updated node status to online.")

}

// updateStatus updates the node status to online
func updateStatus() {
	res, err := util.PostRequest("/node/status/online", fiber.Map{
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
}
