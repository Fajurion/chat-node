package send

import (
	"chat-node/bridge"
	"chat-node/pipe"
)

func sendP2P(sender int64, receiver int64, msg []byte, event pipe.Event) error {

	// Send to receiver
	bridge.Send(receiver, msg)

	return nil
}
