package conversation_routes

import (
	"chat-node/caching"
	"chat-node/database"
	"chat-node/database/conversations"
	"chat-node/util"
	"chat-node/util/requests"

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
	if c.BodyParser(&req) != nil {
		return requests.InvalidRequest(c)
	}

	token, err := caching.ValidateToken(req.ID, req.Token)
	if err != nil {
		return requests.InvalidRequest(c)
	}

	// Check if conversation is group
	var conversation conversations.Conversation
	if database.DBConn.Where("id = ?", token.Conversation).Find(&conversation).Error != nil {
		return requests.InvalidRequest(c)
	}

	if conversation.Type != conversations.TypeGroup {
		return requests.FailedRequest(c, "no.group", nil)
	}

	// Check requirements for a new token
	members, err := caching.LoadMembers(token.Conversation)
	if err != nil {
		return requests.FailedRequest(c, "server.error", err)
	}

	if len(members) >= 100 {
		return requests.FailedRequest(c, "limit.reached", nil)
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

	if database.DBConn.Create(&generated).Error != nil {
		return requests.FailedRequest(c, "server.error", nil)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"id":      generated.ID,
		"token":   generated.Token,
	})
}
