package conversation

import (
	"chat-node/database"
	"chat-node/database/conversations"

	"github.com/Fajurion/pipesfiber/wshandler"
)

// Action: conv_mem
// Action to get all conversation members
func getConversationMembers(message wshandler.Message) {

	if message.ValidateForm("id") {
		wshandler.ErrorResponse(message, "invalid")
		return
	}

	// Check if user is in conversation
	if database.DBConn.Where("conversation = ? AND account = ?", message.Data["id"], message.Client.ID).Find(&conversations.Member{}).Error != nil {
		wshandler.ErrorResponse(message, "invalid")
		return
	}

	// Get conversation members
	var members []conversations.Member
	if database.DBConn.Where("conversation = ?", message.Data["id"]).Find(&members).Error != nil {
		wshandler.ErrorResponse(message, "invalid")
		return
	}

	// Send response
	wshandler.NormalResponse(message, map[string]interface{}{
		"success": true,
		"members": members,
	})
}
