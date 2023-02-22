package receive

import (
	"chat-node/bridge"
	"chat-node/bridge/conversation"
	"chat-node/pipe"
)

func receiveProject(message pipe.Message, msg []byte) {

	pj, err := conversation.GetProject(message.Channel.Target[0])
	if err != nil {
		return
	}

	// Send to receiver
	for member := range pj.Members {
		if member != message.Event.Sender {
			bridge.Send(member, msg)
		}
	}
}
