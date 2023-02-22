package receive

import (
	"chat-node/bridge"
	"chat-node/pipe"
	"log"
)

func receiveBroadcast(message pipe.Message, msg []byte) {

	if message.Event.Name == "ping" {
		log.Println("Received ping from node", message.Event.Data["node"])
		return
	}

	// Send to all receivers
	for _, tg := range message.Channel.Target {

		bridge.Send(tg, msg)
	}
}
