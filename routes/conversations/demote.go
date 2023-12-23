package conversation_routes

import (
	"chat-node/caching"
	"chat-node/database"
	"chat-node/database/conversations"
	message_routes "chat-node/routes/conversations/message"
	"fmt"

	integration "fajurion.com/node-integration"
	"github.com/gofiber/fiber/v2"
)

// Route: /conversations/demote_token
func demoteToken(c *fiber.Ctx) error {

	var req promoteTokenRequest
	if integration.BodyParser(c, &req) != nil {
		return integration.InvalidRequest(c, "invalid request")
	}

	token, err := caching.ValidateToken(req.ID, req.Token)
	if err != nil {
		return integration.InvalidRequest(c, fmt.Sprintf("invalid token: %s", err.Error()))
	}

	// Check if conversation is group
	var conversation conversations.Conversation
	if database.DBConn.Where("id = ?", token.Conversation).Find(&conversation).Error != nil {
		return integration.InvalidRequest(c, fmt.Sprintf("couldn't find conversation in database: %s", err.Error()))
	}

	if conversation.Type != conversations.TypeGroup {
		return integration.FailedRequest(c, "no.group", nil)
	}

	if token.Rank == conversations.RankUser {
		return integration.InvalidRequest(c, "user doesn't have the required rank")
	}

	userToken, err := caching.GetToken(req.User)
	if err != nil {
		return integration.InvalidRequest(c, fmt.Sprintf("specified user doesn't exist: %s", err.Error()))
	}

	if userToken.Conversation != token.Conversation {
		return integration.InvalidRequest(c, "conversations don't match")
	}

	// Get rank to promote (check permissions)
	rankToDemote := userToken.Rank - 1
	if userToken.Rank > token.Rank {
		return integration.InvalidRequest(c, "no permission")
	}

	if err := database.DBConn.Model(&conversations.ConversationToken{}).Where("id = ? AND conversation = ?", userToken.ID, userToken.Conversation).Update("rank", rankToDemote).Error; err != nil {
		return integration.FailedRequest(c, "server.error", err)
	}
	prevRank := userToken.Rank
	userToken.Rank = rankToDemote
	err = caching.UpdateToken(userToken)
	if err != nil {
		return integration.FailedRequest(c, "server.error", err)
	}

	err = message_routes.SendSystemMessage(token.Conversation, "group.rank_change", []string{fmt.Sprintf("%d", prevRank), fmt.Sprintf("%d", userToken.Rank),
		message_routes.AttachAccount(userToken.Data), message_routes.AttachAccount(token.Data)})
	if err != nil {
		return integration.FailedRequest(c, "server.error", err)
	}

	return integration.SuccessfulRequest(c)
}
