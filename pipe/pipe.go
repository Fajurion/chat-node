package pipe

import (
	"chat-node/bridge"
	"chat-node/bridge/conversation"
)

type Event struct {
	Sender  int64                  `json:"sender"`
	Project int64                  `json:"project"`
	Name    string                 `json:"name"`
	Data    map[string]interface{} `json:"data"`
}

func Send(msg []byte, event Event) error {

	// Send to own client
	bridge.Send(event.Sender, msg)

	// TODO: Send to other node if needed

	// Send to project members
	for member := range conversation.GetProject(event.Project).Members {
		if member != event.Sender {

			// Send to node
			bridge.Send(member, msg)
		}
	}

	return nil
}

func Create() {
	// Initalize the pipe and connect to the other nodes
}
