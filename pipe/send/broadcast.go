package send

import (
	"chat-node/pipe"
	"context"

	"nhooyr.io/websocket"
)

func sendBroadcast(message pipe.Message, msg []byte) error {

	// Send to other nodes
	pipe.IterateConnections(func(_ int64, node *websocket.Conn) bool {
		node.Write(context.Background(), websocket.MessageText, msg)
		return true
	})

	return nil
}
