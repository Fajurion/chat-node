package receive

import (
	"chat-node/bridge"
	"chat-node/pipe"
	"chat-node/pipe/receive/processors"
)

func receiveConversation(message pipe.Message) {

	// Send to receivers
	for _, member := range message.Channel.Target {
		if member != message.Event.Sender {

			// Process the message
			msg := processors.ProcessMarshal(&message, member)
			if msg == nil {
				continue
			}

			bridge.Send(member, msg)
		}
	}
}
