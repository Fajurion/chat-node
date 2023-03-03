package receive

import (
	"chat-node/bridge"
	"chat-node/pipe"
)

func receiveConversation(message pipe.Message, msg []byte) {

	// Send to receivers
	for _, member := range message.Channel.Target {
		if member != message.Event.Sender {
			bridge.Send(member, msg)
		}
	}
}
