package space

import (
	"chat-node/caching"

	integration "fajurion.com/node-integration"
	"github.com/Fajurion/pipesfiber/wshandler"
)

// Action: spc_start
func start(message wshandler.Message) {

	if caching.IsInSpace(message.Client.ID) {
		wshandler.ErrorResponse(message, "already.in.space")
		return
	}

	// Create space
	appToken, valid := caching.CreateSpace(message.Client.ID, integration.ClusterID)
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
