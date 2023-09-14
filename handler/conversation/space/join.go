package space

import (
	"chat-node/caching"

	integration "fajurion.com/node-integration"
	"github.com/Fajurion/pipesfiber/wshandler"
)

// Action: cll_join
func joinCall(message wshandler.Message) {

	if message.ValidateForm("id") {
		wshandler.ErrorResponse(message, "invalid")
		return
	}

	if caching.IsInSpace(message.Client.ID) {
		wshandler.ErrorResponse(message, "already.in.space")
		return
	}

	// Create space
	appToken, valid := caching.JoinSpace(message.Client.ID, message.Data["id"].(string), integration.ClusterID)
	if !valid {
		wshandler.ErrorResponse(message, "server.error")
		return
	}

	// Send space info
	wshandler.NormalResponse(message, map[string]interface{}{
		"success": true,
		"token":   appToken,
	})
}
