package receive

import (
	"chat-node/bridge"
	"chat-node/pipe"
	"chat-node/pipe/receive/processors"
	"log"
)

func receiveBroadcast(message pipe.Message) {

	if message.Event.Name == "ping" {
		log.Println("Received ping from node", message.Event.Data["node"])
	}

	// Send to all receivers
	for _, tg := range message.Channel.Target {

		// Process the message
		msg := processors.ProcessMarshal(&message, tg)
		if msg == nil {
			continue
		}

		bridge.Send(tg, msg)
	}
}
