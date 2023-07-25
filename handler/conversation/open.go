package conversation

import (
	"chat-node/database"
	"chat-node/database/conversations"
	"chat-node/database/fetching"
	"chat-node/util"
	"fmt"
	"log"

	"github.com/Fajurion/pipes"
	"github.com/Fajurion/pipes/send"
	"github.com/Fajurion/pipesfiber/wshandler"
	"github.com/bytedance/sonic"
)

// Action: conv_open
func openConversation(message wshandler.Message) {

	if message.ValidateForm("tokens", "data", "keys") {
		wshandler.ErrorResponse(message, "invalid")
		return
	}

	// TODO: Fix conversation opening

	tokens := uint(message.Data["tokens"].(float64))

	if tokens > 100 {
		wshandler.ErrorResponse(message, "member.limit")
		return
	}

	data := message.Data["data"].(string)
	var keys map[string]string = make(map[string]string)

	err := sonic.UnmarshalString(message.Data["keys"].(string), &keys)
	if err != nil {
		wshandler.ErrorResponse(message, "sonic.slipped")
		return
	}

	if len(keys) != len(members)+1 {
		log.Println("keys", len(keys), "members", len(members))
		wshandler.ErrorResponse(message, "invalid")
		return
	}

	// Check if all users are friends
	res, err := util.PostRequest("/account/friends/check", map[string]interface{}{
		"id":      util.NODE_ID,
		"token":   util.NODE_TOKEN,
		"account": message.Client.ID,
		"users":   members,
	})

	if err != nil {
		wshandler.ErrorResponse(message, "server.error")
		return
	}

	if !res["success"].(bool) {

		log.Println("server")

		wshandler.ErrorResponse(message, res["error"].(string))
		return
	}

	members = append(members, message.Client.ID)

	var conversation = conversations.Conversation{
		ID:   util.GenerateToken(12),
		Data: data,
	}

	if database.DBConn.Create(&conversation).Error != nil {
		wshandler.ErrorResponse(message, "server.error")
		return
	}

	var memberList []conversations.Member
	for _, member := range members {

		var role uint = conversations.RoleMember
		if member == message.Client.ID {
			role = conversations.RoleOwner
		}

		var memberObj = conversations.Member{
			ID:           util.GenerateToken(12),
			Conversation: conversation.ID,
			Role:         role,
			Account:      member,
		}
		if database.DBConn.Create(&memberObj).Error != nil {
			wshandler.ErrorResponse(message, "server.error")
			return
		}

		memberList = append(memberList, memberObj)

		// Add stored action
		database.DBConn.Create(&fetching.Action{
			ID:      util.GenerateToken(32),
			Account: member,
			Action:  "conv_key",
			Target:  fmt.Sprintf("%s:%s", conversation.ID, keys[member]),
		})
	}

	// Let the user know that they have a new conversation
	err = send.Pipe(send.ProtocolWS, pipes.Message{
		Channel: pipes.BroadcastChannel(members),
		Event: pipes.Event{
			Name:   "conv_open:l",
			Sender: message.Client.ID,
			Data: map[string]interface{}{
				"success":      true,
				"conversation": conversation,
				"members":      memberList,
				"keys":         keys,
			},
		},
	})

	if err != nil {
		log.Println(err)
		wshandler.ErrorResponse(message, "server.error")
		return
	}

	wshandler.SuccessResponse(message)
}
