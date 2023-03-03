package receive

import (
	"chat-node/pipe"
	"log"

	"github.com/bytedance/sonic"
)

func Handle(message pipe.Message) {

	// Marshal the event
	msg, err := sonic.Marshal(message.Event)
	if err != nil {
		return
	}

	log.Printf("%s: %d: %s", message.Channel.Channel, message.Event.Sender, message.Event.Name)

	switch message.Channel.Channel {
	case "broadcast":
		receiveBroadcast(message, msg)

	case "conversation":
		receiveConversation(message, msg)

	case "p2p":
		receiveP2P(message, msg)

	}
}
