package space

import (
	"chat-node/caching"

	"github.com/Fajurion/pipesfiber/wshandler"
)

// Action: spc_leave
func leaveCall(message wshandler.Message) {

	// Check if in space
	if !caching.IsInSpace(message.Client.ID) {
		wshandler.ErrorResponse(message, "not.in.space")
		return
	}

	// Leave space
	valid := caching.LeaveSpace(message.Client.ID)
	if !valid {
		wshandler.ErrorResponse(message, "server.error")
		return
	}

	// Send success
	wshandler.SuccessResponse(message)
}
