package service

import (
	"chat-node/bridge"
	"chat-node/database"
	"chat-node/database/conversations"
	"chat-node/database/fetching"
	"log"
	"time"

	"github.com/Fajurion/pipes"
)

func setup_mes(client *bridge.Client, current *fetching.Session, account *string) bool {

	// Delete old messages (for now 1 week old messages)
	database.DBConn.
		Where("creation < ? AND EXISTS ( SELECT conversation FROM members AS mem1 WHERE account = ? AND mem1.conversation = messages.conversation )",
			time.Now().Add(-time.Hour*24*7).UnixMilli(), account).
		Delete(&conversations.Message{})

	// Check if the user has any new messages
	var messageList []conversations.Message
	database.DBConn.
		Raw("SELECT * FROM messages AS ms1 WHERE creation > ? AND EXISTS ( SELECT conversation FROM members AS mem1 WHERE account = ? AND mem1.conversation = ms1.conversation )",
			current.LastFetch, account).
		Scan(&messageList)

	log.Println("Messages:", messageList)

	// Send the messages to the user
	client.SendEvent(pipes.Event{
		Name: "setup_msg",
		Data: map[string]interface{}{
			"messages": messageList,
		},
	})

	return true
}
