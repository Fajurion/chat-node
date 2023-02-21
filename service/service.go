package service

import (
	"chat-node/bridge"
	"chat-node/database"
	"chat-node/database/conversations"
	"chat-node/database/fetching"
	"chat-node/pipe"
	"time"
)

func User(client *bridge.Client) bool {
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
	if database.DBConn.Where(&fetching.Session{ID: session}).Take(&current).Error != nil {

		// TODO: New device (sync with old device)

		// Save the session
		if database.DBConn.Create(&fetching.Session{
			ID:        session,
			LastFetch: time.Now().UnixMilli(),
		}).Error != nil {
			return false
		}

		client.SendEvent(pipe.Event{
			Name: "setup",
			Data: map[string]interface{}{
				"message": "new.device",
				"device":  client.Session,
			},
		})
		return true
	}

	// Existing device
	var latest fetching.Session
	if err := database.DBConn.Where(&fetching.Session{
		ID: session,
	}).Order("last_fetch DESC").Take(&latest).Error; err != nil {
		latest = current
	}

	// Check if the user has any new messages
	var conversationList []uint
	if database.DBConn.Model(&conversations.Member{}).Select("conversation").Where("account = ?", account).Find(&conversationList).Error != nil {
		return false
	}

	var messageList []conversations.Message
	if database.DBConn.Where("conversation IN ?", conversationList).Where("creation > ?", current.LastFetch).Find(&messageList).Error != nil {
		return false
	}

	// Save the session
	latest.LastFetch = time.Now().UnixMilli()
	if database.DBConn.Model(&latest).Update("last_fetch", latest.LastFetch).Error != nil {
		bridge.Remove(client.ID, client.Session)
		return false
	}

	// Send the messages to the user
	client.SendEvent(pipe.Event{
		Name: "messages",
		Data: map[string]interface{}{
			"messages": messageList,
		},
	})

	return true
}
