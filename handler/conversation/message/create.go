package message

import (
	"chat-node/database"
	"chat-node/database/conversations"
	"chat-node/handler"
	"chat-node/util"
	"chat-node/util/requests"
	"log"

	"github.com/Fajurion/pipes"
	"github.com/Fajurion/pipes/send"
)

// Action: conv_msg_create
func createMessage(message handler.Message) {

	if message.ValidateForm("conversation", "data") {
		handler.ErrorResponse(message, "invalid")
		return
	}

	conversationId := message.Data["conversation"].(string)
	data := message.Data["data"].(string)
	id := util.GenerateToken(32)

	if conversations.CheckSize(data) {
		handler.ErrorResponse(message, "too.big")
		return
	}

	var conversation conversations.Conversation
	if database.DBConn.Raw("SELECT * FROM conversations AS c1 WHERE EXISTS ( SELECT * FROM members AS m1 WHERE m1.conversation = c1.id AND m1.account = ? ) AND c1.id = ?", message.Client.ID, conversationId).Scan(&conversation).Error != nil {
		handler.ErrorResponse(message, "not.found")
		return
	}

	var stored conversations.Message
	if database.DBConn.Where(&conversations.Message{ID: id, Conversation: conversation.ID}).Take(&stored).Error == nil {

		handler.ErrorResponse(message, "already.exists")
		return
	}

	certificate, err := conversations.GenerateCertificate(id, message.Client.ID)
	if err != nil {
		handler.ErrorResponse(message, "server.error")
		return
	}

	store := conversations.Message{
		ID:           id,
		Conversation: conversation.ID,
		Certificate:  certificate,
		Data:         data,
		Sender:       message.Client.ID,
		Edited:       false,
	}

	if database.DBConn.Create(&store).Error != nil {
		handler.ErrorResponse(message, "server.error")
		return
	}

	// Send to the conversation
	members, nodes, err := requests.LoadConversationDetails(conversation.ID)
	if err != nil {
		handler.ErrorResponse(message, "server.error")
		return
	}

	var event = pipes.Event{
		Sender: "0", // System, to prevent sending to the client twice
		Name:   "conv_msg",
		Data: map[string]interface{}{
			"message": store,
		},
	}

	log.Println("sending..")

	send.Pipe(send.ProtocolWS, pipes.Message{
		Channel: pipes.Conversation(members, nodes),
		Event:   event,
	})

	handler.NormalResponse(message, map[string]interface{}{
		"success": true,
		"id":      id,
	})
}
