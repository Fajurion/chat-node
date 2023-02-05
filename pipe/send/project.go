package send

import (
	"chat-node/bridge"
	"chat-node/bridge/conversation"
	"chat-node/pipe"
)

func sendToProject(project int64, sender int64, msg []byte, event pipe.Event) error {
	for member := range conversation.GetProject(project).Members {
		if member != sender {

			// Send to member
			bridge.Send(member, msg)
		}
	}

	return nil
}
