package main

import (
	"chat-node/pipe"
	"chat-node/routes"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
)

func main() {

	// Create fiber app
	app := fiber.New(fiber.Config{
		JSONEncoder: sonic.Marshal,
		JSONDecoder: sonic.Unmarshal,
	})

	app.Route("/", routes.Setup)

	pipe.Create()

	// Start fiber
	app.Listen(":3001")

}
