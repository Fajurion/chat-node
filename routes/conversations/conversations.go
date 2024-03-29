package conversation_routes

import (
	"chat-node/database"
	"chat-node/database/conversations"
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
	router.Post("/generate_token", generateToken)
	router.Post("/kick_member", kickMember)
	router.Post("/leave", leaveConversation)

	router.Route("/message", message_routes.SetupRoutes)
}

// This deletes all data related to a conversation
func deleteConversation(id string) error {
	if err := database.DBConn.Where("conversation = ?", id).Delete(&conversations.Message{}).Error; err != nil {
		return err
	}
	if err := database.DBConn.Where("conversation = ?", id).Delete(&conversations.ConversationToken{}).Error; err != nil {
		return err
	}
	if err := database.DBConn.Where("id = ?", id).Delete(&conversations.Conversation{}).Error; err != nil {
		return err
	}
	return nil
}
