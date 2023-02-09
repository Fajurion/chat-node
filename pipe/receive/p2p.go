package receive

import (
	"chat-node/bridge"
	"chat-node/pipe"
)

func receiveP2P(message pipe.Message, msg []byte) {

	// Send to receiver
	bridge.Send(message.Channel.Target[0], msg)
}
