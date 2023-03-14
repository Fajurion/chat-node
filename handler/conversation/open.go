package conversation

import (
	"chat-node/database"
	"chat-node/database/conversations"
	"chat-node/handler"
	"chat-node/pipe"
	"chat-node/pipe/send"
	"chat-node/util"
	"fmt"
	"log"
	"time"
)

// Action: conv_open
func openConversation(message handler.Message) {

	if message.ValidateForm("members", "data") {
		handler.ErrorResponse(message, "invalid")
		return
	}

	var members []int64
	for _, member := range message.Data["members"].([]interface{}) {
		members = append(members, int64(member.(float64)))
	}

	data := message.Data["data"].(string)

	log.Println(members)

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

		log.Println("server")

		handler.ErrorResponse(message, res["error"].(string))
		return
	}

	// Enforce limit of 10 conversations per user
	members = append(members, message.Client.ID)

	var conversationCount int64
	if err := database.DBConn.Raw("SELECT COUNT(*) FROM conversations AS c1 WHERE EXISTS ( SELECT id FROM members AS mem1 WHERE conversation = c1.id AND mem1.id IN ? )", members).Scan(&conversationCount).Error; err != nil {

		log.Println(err.Error())

		handler.ErrorResponse(message, "server.error")
		return
	}
	log.Printf("conversation count: %d", conversationCount)

	if conversationCount >= 10 {
		handler.ErrorResponse(message, fmt.Sprintf("limit.reached.%d", conversationCount))
		return
	}

	var conversation = conversations.Conversation{
		Creator:   message.Client.ID,
		Data:      data,
		CreatedAt: time.Now().UnixMilli(),
	}

	if database.DBConn.Create(&conversation).Error != nil {
		handler.ErrorResponse(message, "server.error")
		return
	}

	for _, member := range members {

		var role uint = conversations.RoleMember
		if member == message.Client.ID {
			role = conversations.RoleOwner
		}

		if database.DBConn.Create(&conversations.Member{
			Conversation: conversation.ID,
			Role:         role,
			Account:      member,
		}).Error != nil {
			handler.ErrorResponse(message, "server.error")
			return
		}

	}

	// Let the user know that they have a new conversation
	send.Pipe(pipe.Message{
		Channel: pipe.BroadcastChannel(members),
		Event: pipe.Event{
			Name: "conv_open:l",
			Data: map[string]interface{}{
				"success":      true,
				"conversation": conversation,
				"members":      members,
			},
		},
	})

	handler.SyncSession(message)
	handler.SuccessResponse(message)
}
