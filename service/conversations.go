package service

import (
	"chat-node/database"
	"chat-node/database/conversations"
	"chat-node/database/fetching"

	"github.com/Fajurion/pipes"
	"github.com/Fajurion/pipesfiber"
)

func setup_conv(client *pipesfiber.Client, account *string, current *fetching.Session) bool {

	var conversationList []conversations.Conversation
	database.DBConn.Raw("SELECT * FROM conversations AS c1 WHERE created_at > ? AND EXISTS ( SELECT conversation FROM members AS mem1 WHERE account = ? AND mem1.conversation = c1.id )", current.LastFetch, *account).Scan(&conversationList)

	// Send the conversations to the user
	client.SendEvent(pipes.Event{
		Name: "setup_conv",
		Data: map[string]interface{}{
			"conversations": conversationList,
		},
	})

	// Get members of the conversations
	for _, conversation := range conversationList {
		var memberList []conversations.Member
		if database.DBConn.Where("conversation = ?", conversation.ID).Find(&memberList).Error != nil {
			return false
		}

		// Send the members to the user
		client.SendEvent(pipes.Event{
			Name: "setup_mem",
			Data: map[string]interface{}{
				"conversation": conversation.ID,
				"members":      memberList,
			},
		})
	}

	return true
}
