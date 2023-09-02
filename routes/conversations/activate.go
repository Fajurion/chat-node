package conversation_routes

import (
	"chat-node/database"
	"chat-node/database/conversations"
	"chat-node/util"

	integration "fajurion.com/node-integration"
	"github.com/gofiber/fiber/v2"
)

// Public so it can be unit tested (in the future ig)
type ActivateConversationRequest struct {
	ID    string `json:"id"`
	Token string `json:"token"`
}

func (r *ActivateConversationRequest) Validate() bool {
	return len(r.ID) > 0 && len(r.Token) > 0 && len(r.Token) != util.ConversationTokenLength
}

type returnableMember struct {
	Rank uint   `json:"rank"`
	Data string `json:"data"`
}

// Route: /conversations/activate
func activate(c *fiber.Ctx) error {

	var req ActivateConversationRequest
	if err := c.BodyParser(&req); err != nil {
		return integration.InvalidRequest(c)
	}

	// Validate request
	if !req.Validate() {
		return integration.InvalidRequest(c)
	}

	// Activate conversation
	var token conversations.ConversationToken
	if err := database.DBConn.Where("id = ? AND token = ?", req.ID, req.Token).First(&token).Error; err != nil {
		return integration.InvalidRequest(c)
	}

	if token.Activated {
		return integration.InvalidRequest(c)
	}

	// Activate token
	token.Activated = true
	token.Token = util.GenerateToken(util.ConversationTokenLength)

	if err := database.DBConn.Save(&token).Error; err != nil {
		return integration.FailedRequest(c, "server.error", err)
	}

	// Return all data
	var tokens []conversations.ConversationToken
	if err := database.DBConn.Where(&conversations.ConversationToken{Conversation: token.Conversation}).Find(&tokens).Error; err != nil {
		return integration.FailedRequest(c, "server.error", err)
	}

	var members []returnableMember
	for _, token := range tokens {
		members = append(members, returnableMember{
			Rank: token.Rank,
			Data: token.Data,
		})
	}

	var data string
	if err := database.DBConn.Select("data").Model(&conversations.Conversation{}).Where("id = ?", token.Conversation).Take(&data).Error; err != nil {
		return integration.FailedRequest(c, "server.error", err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    data,
		"members": members,
	})
}
