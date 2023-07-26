package conversation_routes

import "github.com/gofiber/fiber/v2"

func SetupRoutes(router fiber.Router) {
	router.Post("/open", openConversation)
}
