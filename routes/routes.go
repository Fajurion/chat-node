package routes

import (
	"chat-node/routes/auth"
	"chat-node/routes/gateway"
	"chat-node/routes/ping"

	"github.com/gofiber/fiber/v2"
)

func Setup(router fiber.Router) {
	router.Route("/auth", auth.Setup)
	router.Route("/gateway", gateway.SetupRoutes)
	router.Post("/ping", ping.Pong)
}
