package conversation

import (
	"chat-node/database"
	"chat-node/database/conversations"
	"chat-node/handler"
	"chat-node/util"
	"fmt"
)

// Action: conv_open
func openConversation(message handler.Message) {

	if message.ValidateForm("user", "members", "data") {
		handler.ErrorResponse(message, "invalid")
		return
	}

	members := message.Data["members"].([]float64)
	data := message.Data["data"].(string)

	// Check if all users are friends
	res, err := util.PostRequest("/account/friends/check", map[string]interface{}{
		"id":      util.NODE_ID,
		"token":   util.NODE_TOKEN,
		"account": message.Client.ID,
		"users":   members,
	})

	if err != nil {
		handler.ErrorResponse(message, "server.error")
		return
	}

	if !res["success"].(bool) {
		handler.ErrorResponse(message, res["error"].(string))
		return
	}

	// Enforce limit of 10 conversations per group of users
	slice := members[:0]
	slice = append(slice, float64(message.Client.ID))

	var conversationCount int64
	if database.DBConn.Raw("SELECT COUNT(id) FROM conversations AS c1 WHERE EXISTS ( SELECT id FROM members AS mem1 WHERE conversation = c1.id AND mem1.id IN ? )", slice).Scan(&conversationCount).Error != nil {
		handler.ErrorResponse(message, "server.error")
		return
	}

	if conversationCount >= 10 {
		handler.ErrorResponse(message, fmt.Sprintf("limit.reached.%d", conversationCount))
		return
	}

	var conversation = conversations.Conversation{
		Creator: message.Client.ID,
		Data:    data,
	}
	if database.DBConn.Create(&conversation).Error != nil {
		handler.ErrorResponse(message, "server.error")
		return
	}

	for _, member := range members {

		var role uint = conversations.RoleMember
		if member == float64(message.Client.ID) {
			role = conversations.RoleOwner
		}

		if database.DBConn.Create(&conversations.Member{
			Conversation: conversation.ID,
			Role:         role,
		}).Error != nil {
			handler.ErrorResponse(message, "server.error")
			return
		}
	}

}
