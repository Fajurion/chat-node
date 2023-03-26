package receive

import (
	"chat-node/pipe"
	"log"
)

func Handle(message pipe.Message) {

	log.Printf("%s: %d: %s", message.Channel.Channel, message.Event.Sender, message.Event.Name)

	switch message.Channel.Channel {
	case "broadcast":
		receiveBroadcast(message)

	case "conversation":
		receiveConversation(message)

	case "p2p":
		receiveP2P(message)

	}
}
