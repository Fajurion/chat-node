package space

import (
	"chat-node/caching"
	"chat-node/util/localization"

	integration "fajurion.com/node-integration"
	"github.com/Fajurion/pipesfiber/wshandler"
)

// Action: spc_join
func joinCall(message wshandler.Message) {

	if message.ValidateForm("id") {
		wshandler.ErrorResponse(message, localization.InvalidRequest)
		return
	}

	if caching.IsInSpace(message.Client.ID) {
		wshandler.ErrorResponse(message, "already.in.space")
		return
	}

	// Create space
	appToken, valid := caching.JoinSpace(message.Client.ID, message.Data["id"].(string), integration.ClusterID)
	if !valid {
		wshandler.ErrorResponse(message, localization.ErrorServer)
		return
	}

	// Send space info
	wshandler.NormalResponse(message, map[string]interface{}{
		"success": true,
		"token":   appToken,
	})
}
