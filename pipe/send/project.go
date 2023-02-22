package send

import (
	"chat-node/bridge/conversation"
	"chat-node/pipe"
	"context"

	"nhooyr.io/websocket"
)

func sendToProject(message pipe.Message, msg []byte) error {

	pj, err := conversation.GetProject(message.Channel.Target[0])
	if err != nil {
		return err
	}

	for member, node := range pj.Members {
		if member != message.Event.Sender {

			// Send to member
			pipe.GetConnection(node).Write(context.Background(), websocket.MessageText, msg)
		}

	}

	return nil
}
