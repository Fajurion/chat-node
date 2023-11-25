package conversation_routes

import (
	"chat-node/caching"
	"chat-node/util/requests"

	"github.com/gofiber/fiber/v2"
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

	members, err := caching.LoadMembers(token.Conversation)
	if err != nil {
		return requests.InvalidRequest(c)
	}

	realMembers := make([]returnableMemberToken, len(members))
	for i, memberToken := range members {

		member, err := caching.GetToken(memberToken.TokenID)
		if err != nil {
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
