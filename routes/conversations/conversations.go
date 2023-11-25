package conversation_routes

import (
	message_routes "chat-node/routes/conversations/message"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(router fiber.Router) {
	router.Post("/open", openConversation)
	router.Post("/read", read)
	router.Post("/activate", activate)
	router.Post("/demote_token", demoteToken)
	router.Post("/promote_token", promoteToken)
	router.Post("/tokens", listTokens)

	router.Route("/message", message_routes.SetupRoutes)
}
