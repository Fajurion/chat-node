package send

import (
	"chat-node/bridge"
	"chat-node/pipe"
)

func sendBroadcast(receivers []int64, msg []byte, event pipe.Event) error {
	for _, receiver := range receivers {

		// Send to receiver
		bridge.Send(receiver, msg)
	}

	return nil
}
