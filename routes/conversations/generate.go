package conversation_routes

import (
	"chat-node/caching"
	"chat-node/database"
	"chat-node/database/conversations"
	message_routes "chat-node/routes/conversations/message"
	"chat-node/util"
	"fmt"

	integration "fajurion.com/node-integration"
	"github.com/gofiber/fiber/v2"
)

type generateTokenRequest struct {
	ID    string `json:"id"`
	Token string `json:"token"`
	Data  string `json:"data"`
}

// Route: /conversations/generate_token
func generateToken(c *fiber.Ctx) error {

	var req generateTokenRequest
	if integration.BodyParser(c, &req) != nil {
		return integration.InvalidRequest(c, "invalid request")
	}

	token, err := caching.ValidateToken(req.ID, req.Token)
	if err != nil {
		return integration.InvalidRequest(c, fmt.Sprintf("invalid token: %s", err.Error()))
	}

	// Check if conversation is group
	var conversation conversations.Conversation
	if err := database.DBConn.Where("id = ?", token.Conversation).Find(&conversation).Error; err != nil {
		return integration.InvalidRequest(c, fmt.Sprintf("couldn't find conversation in database: %s", err.Error()))
	}

	if conversation.Type != conversations.TypeGroup {
		return integration.FailedRequest(c, "no.group", nil)
	}

	// Check requirements for a new token
	members, err := caching.LoadMembers(token.Conversation)
	if err != nil {
		return integration.FailedRequest(c, "server.error", err)
	}

	if len(members) >= 100 {
		return integration.FailedRequest(c, "limit.reached", nil)
	}

	// Generate a new token
	generated := conversations.ConversationToken{
		ID:           util.GenerateToken(util.ConversationTokenIDLength),
		Token:        util.GenerateToken(util.ConversationTokenLength),
		Activated:    true,
		Conversation: token.Conversation,
		Rank:         conversations.RankAdmin,
		Data:         req.Data,
	}

	if err := database.DBConn.Create(&generated).Error; err != nil {
		return integration.FailedRequest(c, "server.error", err)
	}

	err = message_routes.SendSystemMessage(token.Conversation, message_routes.GroupMemberJoin, []string{message_routes.AttachAccount(token.Data)})
	if err != nil {
		return integration.FailedRequest(c, "server.error", err)
	}

	return integration.ReturnJSON(c, fiber.Map{
		"success": true,
		"id":      generated.ID,
		"token":   generated.Token,
	})
}
