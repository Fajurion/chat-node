package conversation_routes

import (
	"chat-node/caching"
	"chat-node/database"
	"chat-node/database/conversations"
	"chat-node/util/requests"

	"github.com/gofiber/fiber/v2"
)

// Route: /conversations/demote_token
func demoteToken(c *fiber.Ctx) error {

	var req promoteTokenRequest
	if c.BodyParser(&req) != nil {
		return requests.InvalidRequest(c)
	}

	token, err := caching.ValidateToken(req.ID, req.Token)
	if err != nil {
		return requests.InvalidRequest(c)
	}

	if token.Rank == conversations.RankUser {
		return requests.InvalidRequest(c)
	}

	userToken, err := caching.GetToken(req.User)
	if err != nil {
		return requests.InvalidRequest(c)
	}

	// Get rank to promote (check permissions)
	rankToDemote := userToken.Rank - 1
	if userToken.Rank > token.Rank {
		return requests.InvalidRequest(c)
	}

	if database.DBConn.Model(&conversations.ConversationToken{}).Where("id = ? AND conversation = ?", userToken.ID, userToken.Conversation).Update("rank", rankToDemote).Error != nil {
		return requests.FailedRequest(c, "server.error", nil)
	}
	userToken.Rank = rankToDemote
	err = caching.UpdateToken(userToken)
	if err != nil {
		return requests.FailedRequest(c, "server.error", nil)
	}

	return requests.SuccessfulRequest(c)
}
