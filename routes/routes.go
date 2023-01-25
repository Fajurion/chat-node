package routes

import (
	"chat-node/routes/gateway"
	"chat-node/routes/ping"

	"github.com/gofiber/fiber/v2"
)

func Setup(router fiber.Router) {
	router.Route("/gateway", gateway.SetupRoutes)
	router.Post("/ping", ping.Pong)
}
