package receive

import (
	"chat-node/pipe"

	"github.com/bytedance/sonic"
)

func Handle(message pipe.Message) {

	// Marshal the event
	msg, err := sonic.Marshal(message.Event)
	if err != nil {
		return
	}

	switch message.Channel.Channel {
	case "broadcast":
		receiveBroadcast(message, msg)

	case "project":
		receiveProject(message, msg)

	case "p2p":
		receiveP2P(message, msg)

	}
}
