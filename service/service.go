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

	// Check if the account is already in the database
	var status fetching.Status
	if database.DBConn.Where(&fetching.Status{ID: account}).Take(&status).Error != nil {

		// Create a new status
		if database.DBConn.Create(&fetching.Status{
			ID:     account,
			Status: "-",
			Node:   pipe.CurrentNode.ID,
		}).Error != nil {
			return false
		}
	} else {

		// Update the status
		database.DBConn.Model(&fetching.Status{}).Where(&fetching.Status{ID: account}).Update("node", pipe.CurrentNode.ID)
	}

	// Check if this is a new device
	var current fetching.Session
	if database.DBConn.Where(&fetching.Session{ID: session}).Take(&current).Error != nil {

		// TODO: New device (sync with old device)

		// Save the session
		if database.DBConn.Create(&fetching.Session{
			ID:        session,
			Account:   account,
			Node:      pipe.CurrentNode.ID,
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
	if database.DBConn.Model(&latest).Update("last_fetch", latest.LastFetch).Update("node", pipe.CurrentNode.ID).Error != nil {
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
