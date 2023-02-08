package main

import (
	"chat-node/handler"
	"chat-node/pipe"
	"chat-node/routes"
	"chat-node/setup"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {

	// Setting up the node
	if !setup.Setup() {
		return
	}

	// Create fiber app
	app := fiber.New(fiber.Config{
		JSONEncoder: sonic.Marshal,
		JSONDecoder: sonic.Unmarshal,
	})

	app.Use(logger.New())

	app.Route("/", routes.Setup)

	pipe.Create()

	handler.Initialize()

	// Start fiber
	app.Listen("127.0.0.1:" + setup.StartPort)

}
