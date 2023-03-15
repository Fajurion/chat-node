package message

import (
	"chat-node/handler"
	"chat-node/pipe"
	"chat-node/pipe/send"
	"chat-node/util/requests"
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

	if !contains(members, message.Client.ID) {
		return
	}

	send.Pipe(pipe.Message{
		Channel: pipe.Conversation(members, nodes),
		Event: pipe.Event{
			Name:   message.Action,
			Sender: message.Client.ID,
			Data: map[string]interface{}{
				"id": id,
			},
		},
	})
}

func contains(s []int64, e int64) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
