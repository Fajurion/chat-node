package send

import (
	"chat-node/bridge"
	"chat-node/pipe"
	"context"

	"nhooyr.io/websocket"
)

func sendBroadcast(message pipe.Message, msg []byte) error {
	for _, receiver := range message.Channel.Target {

		// Send to receiver
		bridge.Send(receiver, msg)
	}

	// Send to other nodes
	for _, node := range pipe.NodeConnections {
		node.Write(context.Background(), websocket.MessageText, msg)
	}

	return nil
}
