package conversation

import (
	"chat-node/database"
	"chat-node/database/conversations"
	"chat-node/database/fetching"
	"chat-node/handler"
	"chat-node/pipe"
	"chat-node/pipe/send"
	"chat-node/util"
	"fmt"
	"log"
	"time"

	"github.com/bytedance/sonic"
)

// Action: conv_open
func openConversation(message handler.Message) {

	if message.ValidateForm("members", "data", "keys") {
		handler.ErrorResponse(message, "invalid")
		return
	}

	var members []int64
	for _, member := range message.Data["members"].([]interface{}) {
		members = append(members, int64(member.(float64)))
	}

	if len(members) > 100 {
		handler.ErrorResponse(message, "member.limit")
		return
	}

	data := message.Data["data"].(string)
	var keys map[string]string = make(map[string]string)

	err := sonic.UnmarshalString(message.Data["keys"].(string), &keys)
	if err != nil {
		handler.ErrorResponse(message, "sonic.slipped")
		return
	}

	if len(keys) != len(members)+1 {
		log.Println("keys", len(keys), "members", len(members))
		handler.ErrorResponse(message, "invalid")
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
		handler.ErrorResponse(message, "server.error")
		return
	}

	if !res["success"].(bool) {

		log.Println("server")

		handler.ErrorResponse(message, res["error"].(string))
		return
	}

	members = append(members, message.Client.ID)

	var conversation = conversations.Conversation{
		Type:      "chat",
		Creator:   message.Client.ID,
		Data:      data,
		CreatedAt: time.Now().UnixMilli(),
	}

	if database.DBConn.Create(&conversation).Error != nil {
		handler.ErrorResponse(message, "server.error")
		return
	}

	var memberList []conversations.Member
	for _, member := range members {

		var role uint = conversations.RoleMember
		if member == message.Client.ID {
			role = conversations.RoleOwner
		}

		var memberObj = conversations.Member{
			Conversation: conversation.ID,
			Role:         role,
			Account:      member,
		}
		if database.DBConn.Create(&memberObj).Error != nil {
			handler.ErrorResponse(message, "server.error")
			return
		}

		memberList = append(memberList, memberObj)

		// Add stored action
		database.DBConn.Create(&fetching.Action{
			ID:      util.GenerateToken(32),
			Account: member,
			Action:  "conv_key",
			Target:  fmt.Sprintf("%d:%s", conversation.ID, keys[fmt.Sprintf("%d", member)]),
		})
	}

	// Let the user know that they have a new conversation
	send.Pipe(pipe.Message{
		Channel: pipe.BroadcastChannel(members),
		Event: pipe.Event{
			Name: "conv_open:l",
			Data: map[string]interface{}{
				"success":      true,
				"conversation": conversation,
				"members":      memberList,
				"keys":         keys,
			},
		},
	})

	handler.SuccessResponse(message)
}
