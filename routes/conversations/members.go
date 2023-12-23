package conversation_routes

import (
	"chat-node/caching"
	"fmt"

	integration "fajurion.com/node-integration"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type listMembersRequest struct {
	ID    string `json:"id"`
	Token string `json:"token"`
}

type returnableMemberToken struct {
	ID   string `json:"id"`   // Conversation token id
	Data string `json:"data"` // Account id (encrypted)
	Rank uint   `json:"rank"` // Conversation rank
}

// Route: /conversations/tokens
func listTokens(c *fiber.Ctx) error {

	var req listMembersRequest
	if integration.BodyParser(c, &req) != nil {
		return integration.InvalidRequest(c, "invalid request")
	}

	// Validate token
	token, err := caching.ValidateToken(req.ID, req.Token)
	if err != nil {
		return integration.InvalidRequest(c, fmt.Sprintf("invalid conversation token: %s", err.Error()))
	}

	// We use methods without caching here because if a member leaves on a different node, the cache won't be cleared
	members, err := caching.LoadMembersNew(token.Conversation)
	if err != nil {
		return integration.InvalidRequest(c, fmt.Sprintf("couldn't load members: %s", err.Error()))
	}

	realMembers := make([]returnableMemberToken, len(members))
	for i, memberToken := range members {

		member, err := caching.GetTokenNew(memberToken.TokenID)
		if err != nil && err != gorm.ErrRecordNotFound {
			return integration.FailedRequest(c, "server.error", err)
		}

		realMembers[i] = returnableMemberToken{
			ID:   member.ID,
			Data: member.Data,
			Rank: member.Rank,
		}
	}

	return c.JSON(fiber.Map{
		"success": true,
		"members": realMembers,
	})
}
