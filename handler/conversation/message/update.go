package message

import (
	"chat-node/database"
	"chat-node/database/conversations"

	"github.com/Fajurion/pipes"
	"github.com/Fajurion/pipes/send"
	"github.com/Fajurion/pipesfiber/wshandler"
)

// Action: conv_msg_update
func updateMessage(message wshandler.Message) {

	if message.ValidateForm("conversation", "data", "certificate", "id") {
		wshandler.ErrorResponse(message, "invalid")
		return
	}

	conversationId := uint(message.Data["conversation"].(float64))
	data := message.Data["data"].(string)
	certificate := message.Data["certificate"].(string)
	id := message.Data["id"].(string)

	if conversations.CheckSize(data) {
		wshandler.ErrorResponse(message, "too.big")
		return
	}

	var chatMessage conversations.Message
	var conversation conversations.Conversation
	if database.DBConn.Raw("SELECT * FROM conversations AS c1 WHERE EXISTS ( SELECT * FROM members AS m1 WHERE m1.conversation = c1.id AND m1.account = ? ) AND c1.id = ?", message.Client.ID, conversationId).Scan(&conversation).Error != nil {
		wshandler.ErrorResponse(message, "not.found")
		return
	}

	claims, valid := conversations.GetCertificateClaims(certificate)
	if !valid {
		wshandler.ErrorResponse(message, "invalid.certificate")
		return
	}

	if !claims.Valid(id, message.Client.ID) {
		wshandler.ErrorResponse(message, "invalid.certificate")
		return
	}

	if database.DBConn.Where(&conversations.Message{ID: id, Conversation: conversation.ID}).Take(&chatMessage).Error != nil {

		chatMessage = conversations.Message{
			ID:           id,
			Conversation: conversation.ID,
			Certificate:  certificate,
			Data:         data,
			Edited:       true,
		}

		if database.DBConn.Create(&chatMessage).Error != nil {
			wshandler.ErrorResponse(message, "server.error")
			return
		}

		wshandler.NormalResponse(message, map[string]interface{}{
			"success": true,
			"id":      id,
		})
		return
	}

	if database.DBConn.Model(&chatMessage).Update("data", data).Update("edited", true).Error != nil {
		wshandler.ErrorResponse(message, "server.error")
		return
	}

	// Send to the conversation
	members, nodes, err := integration.LoadConversationDetails(conversation.ID)
	if err != nil {
		wshandler.ErrorResponse(message, "server.error")
		return
	}

	send.Pipe(send.ProtocolWS, pipes.Message{
		Channel: pipes.Conversation(members, nodes),
		Event: pipes.Event{
			Name: "conv_msg",
			Data: map[string]interface{}{
				"id":           id,
				"conversation": conversation.ID,
				"data":         data,
				"edited":       true,
			},
		},
	})

	wshandler.NormalResponse(message, map[string]interface{}{
		"success": true,
		"id":      id,
	})
}
