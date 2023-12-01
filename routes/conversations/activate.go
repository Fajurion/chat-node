package conversation_routes

import (
	"chat-node/database"
	"chat-node/database/conversations"
	message_routes "chat-node/routes/conversations/message"
	"chat-node/util"
	"log"

	integration "fajurion.com/node-integration"
	"github.com/gofiber/fiber/v2"
)

// Public so it can be unit tested (in the future ig)
type ActivateConversationRequest struct {
	ID    string `json:"id"`
	Token string `json:"token"`
}

func (r *ActivateConversationRequest) Validate() bool {
	return len(r.ID) > 0 && len(r.Token) > 0 && len(r.Token) == util.ConversationTokenLength
}

type returnableMember struct {
	ID   string `json:"id"`
	Rank uint   `json:"rank"`
	Data string `json:"data"`
}

// Route: /conversations/activate
func activate(c *fiber.Ctx) error {

	var req ActivateConversationRequest
	if err := c.BodyParser(&req); err != nil {
		return integration.InvalidRequest(c, err.Error())
	}

	// Validate request
	if !req.Validate() {
		log.Println(len(req.Token))
		return integration.InvalidRequest(c, "request is invalid")
	}

	// Activate conversation
	var token conversations.ConversationToken
	if database.DBConn.Where("id = ? AND token = ?", req.ID, req.Token).First(&token).Error != nil {
		return integration.FailedRequest(c, "invalid.token", nil)
	}

	if token.Activated {
		return integration.FailedRequest(c, "already.active", nil)
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
			ID:   token.ID,
			Rank: token.Rank,
			Data: token.Data,
		})
	}

	var conversation conversations.Conversation
	if err := database.DBConn.Where("id = ?", token.Conversation).Take(&conversation).Error; err != nil {
		return integration.FailedRequest(c, "server.error", err)
	}

	err := message_routes.SendSystemMessage(token.Conversation, "group.member_join", []string{message_routes.AttachAccount(token.Data)})
	if err != nil {
		return integration.FailedRequest(c, "server.error", err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"type":    conversation.Type,
		"data":    conversation.Data,
		"token":   token.Token,
		"members": members,
	})
}
