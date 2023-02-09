package receive

import (
	"chat-node/bridge"
	"chat-node/pipe"
)

func receiveBroadcast(message pipe.Message, msg []byte) {

	// Send to all receivers
	for _, tg := range message.Channel.Target {

		bridge.Send(tg, msg)
	}
}
