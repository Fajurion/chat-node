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
		Name: "setup_wel",
		Data: map[string]interface{}{
			"name": client.Username,
			"tag":  client.Tag,
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
		current = fetching.Session{
			ID:        session,
			Account:   account,
			Node:      pipe.CurrentNode.ID,
			LastFetch: 0,
		}

		if database.DBConn.Create(&current).Error != nil {
			return false
		}

		client.SendEvent(pipe.Event{
			Name: "setup_device",
			Data: map[string]interface{}{
				"device": client.Session,
			},
		})
	}

	// Existing device

	// Get new conversations
	var conversationList []conversations.Conversation
	database.DBConn.Where("created_at > ?", current.LastFetch).Take(&conversationList)

	// Send the conversations to the user
	client.SendEvent(pipe.Event{
		Name: "setup_conv",
		Data: map[string]interface{}{
			"conversations": conversationList,
		},
	})

	// Check if the user has any new messages
	var messageList []conversations.Message
	database.DBConn.Raw("SELECT * FROM messages AS ms1 WHERE creation > ? AND EXISTS ( SELECT conversation FROM members AS mem1 WHERE account = ? AND mem1.conversation = ms1.conversation )", current.LastFetch, account).Scan(&messageList)

	/*
		var conversationList []uint
		if database.DBConn.Model(&conversations.Member{}).Select("conversation").Where("account = ?", account).Find(&conversationList).Error != nil {
			return false
		}

		if database.DBConn.Where("conversation IN ?", conversationList).Where("creation > ?", current.LastFetch).Find(&messageList).Error != nil {
			return false
		}
	*/

	// Get new actions
	var actionList []fetching.Action
	database.DBConn.Where("created_at > ?", current.LastFetch).Take(&actionList)

	// Send the actions to the user
	client.SendEvent(pipe.Event{
		Name: "setup_act",
		Data: map[string]interface{}{
			"actions": actionList,
		},
	})

	// Save the session
	current.LastFetch = time.Now().UnixMilli()
	current.Node = pipe.CurrentNode.ID
	if database.DBConn.Save(&current).Error != nil {
		bridge.Remove(client.ID, client.Session)
		return false
	}

	// Send the messages to the user
	client.SendEvent(pipe.Event{
		Name: "setup_msg",
		Data: map[string]interface{}{
			"messages": messageList,
		},
	})

	return true
}
