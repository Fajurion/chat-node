package account

import (
	"chat-node/database"
	"chat-node/database/fetching"
	"chat-node/handler/conversation"
	"chat-node/util"
	"log"

	"github.com/Fajurion/pipes"
	"github.com/Fajurion/pipes/send"
	"github.com/Fajurion/pipesfiber/wshandler"
)

// Action: st_send
func sendStatus(message wshandler.Message) {

	if message.ValidateForm("tokens", "status", "data") {
		wshandler.ErrorResponse(message, "invalid")
		return
	}

	// Save in database
	statusMessage := message.Data["status"].(string)
	data := message.Data["data"].(string)
	if err := database.DBConn.Model(&fetching.Status{}).Where("id = ?", message.Client.ID).Update("data", statusMessage).Error; err != nil {
		wshandler.ErrorResponse(message, "server.error")
		return
	}

	// Send to other people
	conversationTokens, _, members, _, ok := conversation.PrepareConversationTokens(message)
	if !ok {
		return
	}

	for _, token := range conversationTokens {

		var memberIds []string
		var memberNodes []string
		log.Printf("%d", len(members[token.Conversation]))
		if len(members[token.Conversation]) == 2 {
			for _, member := range members[token.Conversation] {
				if member.Token != token.Token {
					memberIds = append(memberIds, "s-"+member.Token)
					memberNodes = append(memberNodes, util.Node64(member.Node))
				}
			}
		}
		log.Printf("Sending to %d members", len(memberIds))

		// Send the subscription event
		send.Pipe(send.ProtocolWS, pipes.Message{
			Channel: pipes.Conversation(memberIds, memberNodes),
			Event:   statusEvent(statusMessage, data, ""),
		})
	}

	wshandler.SuccessResponse(message)
}

func statusEvent(st string, data string, suffix string) pipes.Event {
	return pipes.Event{
		Name: "acc_st" + suffix,
		Data: map[string]interface{}{
			"st": st,
			"d":  data,
		},
	}
}
