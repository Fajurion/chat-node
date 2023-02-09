package send

import (
	"chat-node/bridge"
	"chat-node/pipe"
)

func sendP2P(message pipe.Message, msg []byte) error {

	// Check if receiver is on this node
	if _, ok := bridge.Connections[message.Channel.Target[0]]; ok {
		bridge.Send(message.Channel.Sender, msg)
		return nil
	}

	// TODO: Get other node

	return nil
}
