package conversation_routes

import (
	"chat-node/database"
	"chat-node/database/conversations"
	"chat-node/util"

	integration "fajurion.com/node-integration"
	"github.com/gofiber/fiber/v2"
)

// Public so it can be unit tested (in the future ig)
type OpenConversationRequest struct {
	AccountData string   `json:"accountData"` // Account data of the user opening the conversation (encrypted)
	Members     []string `json:"members"`
	Data        string   `json:"data"` // Encrypted data
}

func (r *OpenConversationRequest) Validate() bool {
	return len(r.Members) > 0 && len(r.Data) > 0 && len(r.Data) <= util.MaxConversationDataLength
}

type returnableToken struct {
	ID    string `json:"id"`
	Token string `json:"token"`
}

// Route: /conversations/open
func openConversation(c *fiber.Ctx) error {

	var req OpenConversationRequest
	if err := c.BodyParser(&req); err != nil {
		return integration.InvalidRequest(c, err.Error())
	}

	// Validate request
	if !req.Validate() {
		return integration.InvalidRequest(c, "request is invalid")
	}

	if len(req.Members)+1 > util.MaxConversationMembers {
		return integration.FailedRequest(c, "member.limit", nil)
	}

	if len(req.AccountData) > util.MaxConversationTokenDataLength {
		return integration.FailedRequest(c, "data.limit", nil)
	}

	for _, member := range req.Members {
		if len(member) > util.MaxConversationTokenDataLength {
			return integration.FailedRequest(c, "data.limit", nil)
		}
	}

	// Create conversation
	conv := conversations.Conversation{
		ID:   util.GenerateToken(util.ConversationIDLength),
		Data: req.Data,
	}

	if err := database.DBConn.Create(&conv).Error; err != nil {
		return integration.FailedRequest(c, "database.error", nil)
	}

	// Create tokens
	var tokens map[string]returnableToken = make(map[string]returnableToken)
	for _, memberData := range req.Members {

		convToken := util.GenerateToken(util.ConversationTokenLength)

		tk := conversations.ConversationToken{
			ID:           util.GenerateToken(util.ConversationTokenIDLength),
			Conversation: conv.ID,
			Activated:    false,
			Token:        convToken,
			Rank:         conversations.RankUser,
			Data:         memberData,
			LastRead:     0,
		}

		if err := database.DBConn.Create(&tk).Error; err != nil {
			return integration.FailedRequest(c, "server.error", err)
		}

		tokens[util.HashString(memberData)] = returnableToken{
			ID:    tk.ID,
			Token: convToken,
		}
	}

	adminToken := conversations.ConversationToken{
		ID:           util.GenerateToken(util.ConversationTokenIDLength),
		Token:        util.GenerateToken(util.ConversationTokenLength),
		Activated:    true,
		Conversation: conv.ID,
		Rank:         conversations.RankAdmin,
		Data:         req.AccountData,
	}

	if err := database.DBConn.Create(&adminToken).Error; err != nil {
		return integration.FailedRequest(c, "server.error", err)
	}

	// TODO: Fix that the admin can pretend to be one of the users
	return c.JSON(fiber.Map{
		"success":      true,
		"conversation": conv.ID,
		"admin_token": returnableToken{
			ID:    adminToken.ID,
			Token: adminToken.Token,
		},
		"tokens": tokens,
	})
}
