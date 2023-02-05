package main

import (
	"chat-node/handler/setup"
	"chat-node/pipe"
	"chat-node/routes"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {

	// Create fiber app
	app := fiber.New(fiber.Config{
		JSONEncoder: sonic.Marshal,
		JSONDecoder: sonic.Unmarshal,
	})

	app.Use(logger.New())

	app.Route("/", routes.Setup)

	pipe.Create()

	setup.Initialize()

	// Start fiber
	app.Listen("127.0.0.1:3001")

}
