package service

import (
	"chat-node/bridge"
	"chat-node/database"
	"chat-node/database/conversations"
	"chat-node/database/fetching"
	"chat-node/pipe"
)

func User(client *bridge.Client) {
	session := client.Session
	account := client.ID

	client.SendEvent(pipe.Event{
		Name: "setup",
		Data: map[string]interface{}{
			"message":  "welcome",
			"username": client.Username,
			"tag":      client.Tag,
		},
	})

	// Check if this is a new device
	var current fetching.Session
	if err := database.DBConn.Where("session = ?", session).Take(&session).Error; err == nil {

		// TODO: New device

		client.SendEvent(pipe.Event{
			Name: "setup",
			Data: map[string]interface{}{
				"message": "new.device",
				"device":  client.Session,
			},
		})
		return
	}

	// Existing device
	var latest fetching.Session
	if err := database.DBConn.Where("account = ?", account).Order("fetch DESC").Take(&latest).Error; err != nil {
		latest = current
	}

	// Check if the user has any new messages
	var conversationList []uint
	if err := database.DBConn.Model(&conversations.Member{}).Select("conversation").Where("account = ?", account).Find(&conversationList).Error; err != nil {
		return
	}

	var messageList []conversations.Message
	if err := database.DBConn.Where("conversation IN ?", conversationList).Where("creation > ?", current.Fetch).Find(&messageList).Error; err != nil {
		return
	}

	// Send the messages to the user
	client.SendEvent(pipe.Event{
		Name: "messages",
		Data: map[string]interface{}{
			"messages": messageList,
		},
	})
}
