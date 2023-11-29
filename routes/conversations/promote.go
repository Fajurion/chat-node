package conversation_routes

import (
	"chat-node/caching"
	"chat-node/database"
	"chat-node/database/conversations"
	message_routes "chat-node/routes/conversations/message"
	"chat-node/util/requests"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type promoteTokenRequest struct {
	ID    string `json:"id"`
	Token string `json:"token"`
	User  string `json:"user"`
}

// Route: /conversations/promote_token
func promoteToken(c *fiber.Ctx) error {

	var req promoteTokenRequest
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

	if token.Rank == conversations.RankUser {
		return requests.InvalidRequest(c)
	}

	userToken, err := caching.GetToken(req.User)
	if err != nil {
		return requests.InvalidRequest(c)
	}

	if userToken.Conversation != token.Conversation {
		return requests.InvalidRequest(c)
	}

	// Get rank to promote (check permissions)
	rankToPromote := userToken.Rank + 1
	if rankToPromote > token.Rank {
		return requests.InvalidRequest(c)
	}

	if database.DBConn.Model(&conversations.ConversationToken{}).Where("id = ? AND conversation = ?", userToken.ID, userToken.Conversation).Update("rank", rankToPromote).Error != nil {
		return requests.FailedRequest(c, "server.error", nil)
	}
	prevRank := userToken.Rank
	userToken.Rank = rankToPromote
	err = caching.UpdateToken(userToken)
	if err != nil {
		return requests.FailedRequest(c, "server.error", nil)
	}

	err = message_routes.SendSystemMessage(token.Conversation, "group.rank_change", []string{fmt.Sprintf("%d", prevRank), fmt.Sprintf("%d", userToken.Rank), userToken.ID, token.ID})
	if err != nil {
		return requests.FailedRequest(c, "server.error", nil)
	}

	return requests.SuccessfulRequest(c)
}
