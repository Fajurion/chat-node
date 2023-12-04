package conversation_routes

import (
	"chat-node/caching"
	"chat-node/database"
	"chat-node/database/conversations"
	message_routes "chat-node/routes/conversations/message"
	"chat-node/util/requests"

	"github.com/gofiber/fiber/v2"
)

type leaveRequest struct {
	ID    string `json:"id"`
	Token string `json:"token"`
}

// Routes: /conversations/leave
func leaveConversation(c *fiber.Ctx) error {

	var req leaveRequest
	if err := c.BodyParser(&req); err != nil {
		return requests.InvalidRequest(c)
	}

	token, err := caching.ValidateToken(req.ID, req.Token)
	if err != nil {
		return requests.InvalidRequest(c)
	}

	err = message_routes.SendSystemMessage(token.Conversation, "group.member_leave", []string{
		message_routes.AttachAccount(token.Data),
	})
	if err != nil {
		return requests.FailedRequest(c, "server.error", nil)
	}
	caching.DeleteToken(token.ID)

	// Delete token
	if database.DBConn.Where("id = ?", token.ID).Delete(&conversations.ConversationToken{}).Error != nil {
		return requests.FailedRequest(c, "server.error", nil)
	}

	return requests.SuccessfulRequest(c)
}
