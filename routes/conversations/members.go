package conversation_routes

import (
	"chat-node/caching"
	"chat-node/util/requests"

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
	if c.BodyParser(&req) != nil {
		return requests.InvalidRequest(c)
	}

	// Validate token
	token, err := caching.ValidateToken(req.ID, req.Token)
	if err != nil {
		return requests.InvalidRequest(c)
	}

	// We use methods without caching here because if a member leaves on a different node, the cache won't be cleared
	members, err := caching.LoadMembersNew(token.Conversation)
	if err != nil {
		return requests.InvalidRequest(c)
	}

	realMembers := make([]returnableMemberToken, len(members))
	for i, memberToken := range members {

		member, err := caching.GetTokenNew(memberToken.TokenID)
		if err != nil && err != gorm.ErrRecordNotFound {
			return requests.FailedRequest(c, "server.error", err)
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
