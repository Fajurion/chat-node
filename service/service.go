package service

import (
	"chat-node/bridge"
	"chat-node/database"
	"chat-node/database/conversations"
	"chat-node/database/fetching"
	"chat-node/pipe"
	"log"
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

	// Get the earliest fetch time
	var firstFetch int64
	database.DBConn.Raw("SELECT MIN(last_fetch) FROM sessions WHERE account = ?", account).Scan(&firstFetch)

	// Get new conversations
	if !setup_conv(client, &account, &current) {
		return false
	}

	if !setup_act(client, &current, &firstFetch) {
		return false
	}

	// Check if the user has any new messages
	var messageList []conversations.Message
	database.DBConn.Raw("SELECT * FROM messages AS ms1 WHERE creation > ? AND EXISTS ( SELECT conversation FROM members AS mem1 WHERE account = ? AND mem1.conversation = ms1.conversation )", current.LastFetch, account).Scan(&messageList)

	log.Println("Messages:", messageList)

	// Send the messages to the user
	client.SendEvent(pipe.Event{
		Name: "setup_msg",
		Data: map[string]interface{}{
			"messages": messageList,
		},
	})

	// Save the session
	current.LastFetch = time.Now().UnixMilli()
	current.Node = pipe.CurrentNode.ID
	if database.DBConn.Save(&current).Error != nil {
		bridge.Remove(client.ID, client.Session)
		return false
	}

	return true
}
