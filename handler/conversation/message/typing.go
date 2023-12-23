package message

import (
	"github.com/Fajurion/pipes"
	"github.com/Fajurion/pipes/send"
	"github.com/Fajurion/pipesfiber/wshandler"
)

// Action: conv_t_s / conv_t
func typingStatus(message wshandler.Message) {

	if message.ValidateForm("id") {
		return
	}

	id := message.Data["id"].(string)

	// Send to the conversation
	members, nodes, err := integration.LoadConversationDetails(id)
	if err != nil {
		return
	}

	if !contains(members, message.Client.ID) {
		return
	}

	send.Pipe(send.ProtocolWS, pipes.Message{
		Channel: pipes.Conversation(members, nodes),
		Event: pipes.Event{
			Name: message.Action,
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
