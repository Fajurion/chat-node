package send

import (
	"chat-node/bridge"
	"chat-node/pipe"
	"chat-node/util"
	"context"

	"nhooyr.io/websocket"
)

func sendP2P(message pipe.Message, msg []byte) error {

	if _, ok := bridge.Connections.Get(message.Channel.Sender); ok {
		bridge.Send(message.Channel.Sender, msg)
		return nil
	}

	// Check if receiver is on this node
	if message.Channel.Target[0] == int64(util.NODE_ID) {
		bridge.Send(message.Channel.Target[0], msg)
		return nil
	}

	pipe.NodeConnections[message.Channel.Target[1]].Write(context.Background(), websocket.MessageText, msg)

	return nil
}
