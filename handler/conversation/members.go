package conversation

import (
	"chat-node/database"
	"chat-node/database/conversations"
	"chat-node/handler"
)

// Action: conv_mem
// Action to get all conversation members
func getConversationMembers(message handler.Message) {

	if message.ValidateForm("id") {
		handler.ErrorResponse(message, "invalid")
		return
	}

	// Check if user is in conversation
	if database.DBConn.Where("conversation = ? AND account = ?", message.Data["id"], message.Client.ID).Find(&conversations.Member{}).Error != nil {
		handler.ErrorResponse(message, "invalid")
		return
	}

	// Get conversation members
	var members []conversations.Member
	if database.DBConn.Where("conversation = ?", message.Data["id"]).Find(&members).Error != nil {
		handler.ErrorResponse(message, "invalid")
		return
	}

	// Send response
	handler.NormalResponse(message, map[string]interface{}{
		"success": true,
		"members": members,
	})
}
