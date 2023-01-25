package main

import (
	"chat-node/bridge/conversation"
	"chat-node/pipe"
	"chat-node/routes"
	"fmt"

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

	proj := conversation.Project{
		ID:      0,
		Members: make(map[int64]int64),
	}

	fmt.Print(proj.Members[0])

	pipe.Create()

	// Start fiber
	app.Listen(":3000")

}
