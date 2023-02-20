package send

import (
	"chat-node/bridge"
	"chat-node/pipe"

	"github.com/bytedance/sonic"
)

func Pipe(message pipe.Message) error {

	msg, err := sonic.Marshal(message)
	if err != nil {
		return err
	}

	// Send to own client
	if message.Channel.IsProject() {
		bridge.Send(message.Channel.Sender, msg)
	}

	switch message.Channel.Channel {
	case "project":
		return sendToProject(message, msg)

	case "broadcast":
		return sendBroadcast(message, msg)

	case "p2p":
		return sendP2P(message, msg)

	case "client":
		bridge.Send(message.Channel.Target[0], msg)
		return nil
	}

	return nil
}
