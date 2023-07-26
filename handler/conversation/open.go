package conversation

import (
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

	// data := message.Data["data"].(string)
	var keys map[string]string = make(map[string]string)

	err := sonic.UnmarshalString(message.Data["keys"].(string), &keys)
	if err != nil {
		wshandler.ErrorResponse(message, "sonic.slipped")
		return
	}

	// TODO: Create a new conversation

	// Let the user know that they have a new conversation
	/*
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
	*/
}
