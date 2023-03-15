package send

import (
	"chat-node/pipe"
	"context"

	"nhooyr.io/websocket"
)

func sendToConversation(message pipe.Message, msg []byte) error {

	for _, node := range message.Channel.Nodes {
		if node == pipe.CurrentNode.ID {
			continue
		}

		pipe.GetConnection(node).Write(context.Background(), websocket.MessageText, msg)
	}

	return nil
}
