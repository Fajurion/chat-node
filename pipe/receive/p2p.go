package receive

import (
	"chat-node/bridge"
	"chat-node/pipe"
	"chat-node/pipe/receive/processors"
)

func receiveP2P(message pipe.Message) {

	// Process the message
	msg := processors.ProcessMarshal(&message, message.Channel.Target[0])
	if msg == nil {
		return
	}

	// Send to receiver
	bridge.Send(message.Channel.Target[0], msg)
}
