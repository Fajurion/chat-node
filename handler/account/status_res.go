package account

import (
	"chat-node/caching"
	"chat-node/util/localization"

	"github.com/Fajurion/pipes"
	"github.com/Fajurion/pipes/send"
	"github.com/Fajurion/pipesfiber/wshandler"
)

// Action: st_res
func respondToStatus(message wshandler.Message) {

	if message.ValidateForm("id", "token", "status", "data") {
		wshandler.ErrorResponse(message, localization.InvalidRequest)
		return
	}

	id := message.Data["id"].(string)
	token := message.Data["token"].(string)
	status := message.Data["status"].(string)
	data := message.Data["data"].(string)

	// Get from cache
	convToken, err := caching.ValidateToken(id, token)
	if err != nil {
		wshandler.ErrorResponse(message, localization.InvalidRequest)
		return
	}

	// Check if this is a valid conversation
	members, err := caching.LoadMembers(convToken.Conversation)
	if err != nil {
		wshandler.ErrorResponse(message, localization.ErrorServer)
		return
	}

	ids, nodes := caching.MembersToPipes(members)

	// Send the subscription event
	err = send.Pipe(send.ProtocolWS, pipes.Message{
		Channel: pipes.Conversation(ids, nodes),
		Event:   statusEvent(status, data, ":a"),
	})
	if err != nil {
		wshandler.ErrorResponse(message, localization.ErrorServer)
		return
	}

	wshandler.SuccessResponse(message)
}
