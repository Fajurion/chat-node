package pipe

import (
	"chat-node/util"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
)

func Create() {

	log.Println("Creating pipe...")

	// Create pipe

	// Tell server new status
	res, err := util.PostRequest("/node/status/update", fiber.Map{
		"token":  util.NODE_TOKEN,
		"status": util.StatusOnline,
	})

	if err != nil {
		log.Println("Backend is currently offline!")
		os.Exit(1)
	}

	if !res["success"].(bool) {
		log.Println("This node may not be registered..")
		os.Exit(1)
	}

	log.Println("Updated node status to online.")

}
