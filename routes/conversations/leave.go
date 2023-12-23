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

type leaveRequest struct {
	ID    string `json:"id"`
	Token string `json:"token"`
}

// Routes: /conversations/leave
func leaveConversation(c *fiber.Ctx) error {

	var req leaveRequest
	if err := integration.BodyParser(c, &req); err != nil {
		return integration.InvalidRequest(c, "invalid request")
	}

	token, err := caching.ValidateToken(req.ID, req.Token)
	if err != nil {
		return integration.InvalidRequest(c, fmt.Sprintf("invalid conversation token: %s", err.Error()))
	}

	// Delete token
	if err := database.DBConn.Where("id = ?", token.ID).Delete(&conversations.ConversationToken{}).Error; err != nil {
		return integration.FailedRequest(c, "server.error", err)
	}

	if err != nil {
		return integration.FailedRequest(c, "server.error", err)
	}
	caching.DeleteToken(token.ID)

	members, err := caching.LoadMembersNew(token.Conversation)
	if err != nil {
		integration.FailedRequest(c, "server.error", err)
	}

	// Check if the chat is a DM (send delete message if it is)
	var conversation conversations.Conversation
	if err := database.DBConn.Where("id = ?", token.Conversation).Take(&conversation).Error; err != nil {
		return integration.FailedRequest(c, "server.error", err)
	}

	if conversation.Type == conversations.TypePrivateMessage && len(members) == 1 {

		// Send deletion message (this will automatically get rid of the conversation because the other guy will leave after)
		err := message_routes.SendSystemMessage(token.Conversation, "conv.deleted", []string{})
		if err != nil {
			return integration.FailedRequest(c, "server.error", err)
		}

		return integration.SuccessfulRequest(c)
	}

	if len(members) == 0 {

		// Delete conversation
		if err := database.DBConn.Delete(&conversations.Conversation{}, "id = ?", token.Conversation).Error; err != nil {
			return integration.FailedRequest(c, "server.error", err)
		}

		return integration.SuccessfulRequest(c)
	} else {

		// Check if another admin is needed
		if token.Rank == conversations.RankAdmin {
			needed := true
			bestCase := conversations.ConversationToken{
				Rank: conversations.RankUser,
			}
			for _, member := range members {
				userToken, err := caching.GetToken(member.TokenID)
				if err != nil {
					continue
				}

				if userToken.Rank == conversations.RankAdmin {
					needed = false
					break
				}

				if bestCase.Rank <= userToken.Rank {
					bestCase = userToken
				}
			}

			// Promote to admin if needed
			if needed {
				if database.DBConn.Model(&conversations.ConversationToken{}).Where("id = ?", bestCase.ID).Update("rank", conversations.RankAdmin).Error != nil {
					return integration.FailedRequest(c, "server.error", nil)
				}
				err = caching.UpdateToken(bestCase)
				if err != nil {
					return integration.FailedRequest(c, "server.error", nil)
				}

				err = message_routes.SendSystemMessage(token.Conversation, "group.new_admin", []string{message_routes.AttachAccount(bestCase.Data)})
				if err != nil {
					return integration.FailedRequest(c, "server.error", nil)
				}
			}
		}
	}

	message_routes.SendSystemMessage(token.Conversation, "group.member_leave", []string{
		message_routes.AttachAccount(token.Data),
	})

	return integration.SuccessfulRequest(c)
}
