package conversation_routes

import (
	"chat-node/caching"
	message_routes "chat-node/routes/conversations/message"
	"chat-node/util/localization"

	integration "fajurion.com/node-integration"
	"github.com/gofiber/fiber/v2"
)

type kickMemberRequest struct {
	Id     string `json:"id"`
	Token  string `json:"token"`
	Target string `json:"target"`
}

// Route: /conversations/kick_member
func kickMember(c *fiber.Ctx) error {

	var req kickMemberRequest
	if err := integration.BodyParser(c, &req); err != nil {
		return integration.InvalidRequest(c, "invalid request")
	}

	token, err := caching.ValidateToken(req.Id, req.Token)
	if err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}
	targetToken, err := caching.GetToken(req.Target)
	if err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	// Check if the token has the permission
	if token.Rank > targetToken.Rank {
		return integration.FailedRequest(c, localization.KickNoPermission, nil)
	}

	err = message_routes.SendSystemMessage(token.Conversation, message_routes.GroupMemberKick, []string{message_routes.AttachAccount(targetToken.Data)})
	if err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	// TODO: Kick the guy

	return integration.SuccessfulRequest(c)
}
