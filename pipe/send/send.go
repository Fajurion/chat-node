package send

import (
	"chat-node/bridge"
	"chat-node/pipe"
)

func Pipe(channel pipe.Channel, msg []byte, event pipe.Event) error {

	// Send to own client
	if channel.Sender != 0 {
		bridge.Send(channel.Sender, msg)
	}

	switch channel.Channel {
	case "project":
		return sendToProject(channel.Target[0], channel.Sender, msg, event)

	case "broadcast":
		return sendBroadcast(channel.Target, msg, event)

	case "p2p":
		return sendP2P(channel.Sender, channel.Target[0], msg, event)
	}

	return nil
}
