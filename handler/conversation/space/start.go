package space

import (
	"chat-node/caching"
	"chat-node/util/localization"

	integration "fajurion.com/node-integration"
	"github.com/Fajurion/pipesfiber/wshandler"
)

// Action: spc_start
func start(message wshandler.Message) {

	/*
		TODO: Re-enable
		if caching.IsInSpace(message.Client.ID) {
			wshandler.ErrorResponse(message, "already.in.space")
			return
		}
	*/

	// Create space
	roomId, appToken, valid := caching.CreateSpace(message.Client.ID, integration.ClusterID)
	if !valid {
		wshandler.ErrorResponse(message, localization.ErrorServer)
		return
	}

	// Send space info
	wshandler.NormalResponse(message, map[string]interface{}{
		"success": true,
		"id":      roomId,
		"token":   appToken,
	})
}
