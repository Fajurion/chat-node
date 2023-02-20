package send

import (
	"chat-node/bridge"
	"chat-node/pipe"
	"context"

	"nhooyr.io/websocket"
)

func sendBroadcast(message pipe.Message, msg []byte) error {
	for _, receiver := range message.Channel.Target {

		if _, ok := bridge.Connections.Get(receiver); !ok {
			continue
		}

		// Send to receiver
		bridge.Send(receiver, msg)
	}

	// Send to other nodes
	pipe.IterateConnections(func(_ int64, node *websocket.Conn) bool {
		node.Write(context.Background(), websocket.MessageText, msg)
		return true
	})

	return nil
}
