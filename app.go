package main

import (
	"chat-node/database"
	handlerCreate "chat-node/handler/create"
	"chat-node/pipe"
	processors "chat-node/pipe/receive/processors/create"
	"chat-node/routes"
	"chat-node/setup"
	"chat-node/util"
	"log"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {

	// Setting up the node
	if !setup.Setup() {
		return
	}

	// Connect to the database
	database.Connect()

	// Create fiber app
	app := fiber.New(fiber.Config{
		JSONEncoder: sonic.Marshal,
		JSONDecoder: sonic.Unmarshal,
	})

	log.Println(util.GenerateToken(200))

	app.Use(logger.New())

	app.Route("/", routes.Setup)

	pipe.Create()

	// Create handlers
	handlerCreate.Create()

	// Initialize processors
	processors.SetupProcessors()

	// Start fiber
	app.Listen(pipe.CurrentNode.Domain)

}
