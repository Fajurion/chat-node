package message

import (
	"chat-node/handler"
	"chat-node/util"
	"chat-node/util/requests"

	"github.com/Fajurion/pipes"
	"github.com/Fajurion/pipes/send"
)

// Action: conv_t_s / conv_t
func typingStatus(message handler.Message) {

	if message.ValidateForm("id") {
		return
	}

	id := uint(message.Data["id"].(float64))

	// Send to the conversation
	members, nodes, err := requests.LoadConversationDetails(id)
	if err != nil {
		return
	}

	if !contains(members, util.User64(message.Client.ID)) {
		return
	}

	send.Pipe(send.ProtocolWS, pipes.Message{
		Channel: pipes.Conversation(members, nodes),
		Event: pipes.Event{
			Name:   message.Action,
			Sender: util.User64(message.Client.ID),
			Data: map[string]interface{}{
				"id": id,
			},
		},
	})
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
