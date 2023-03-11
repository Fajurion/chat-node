package handler

import (
	"chat-node/bridge"
	"chat-node/pipe"
	"chat-node/pipe/send"
	"log"
	"time"
)

type Message struct {
	Client *bridge.Client         `json:"client"`
	Action string                 `json:"action"` // The action to perform
	Data   map[string]interface{} `json:"data"`
}

// Routes is a map of all the routes
var Routes map[string]func(Message)

func Handle(message Message) bool {
	defer func() {
		if err := recover(); err != nil {
			ErrorResponse(message, "internal")
		}
	}()

	// Check if the action exists
	if Routes[message.Action] == nil {
		return false
	}

	log.Println("Handling message: " + message.Action)

	go Route(message.Action, message)

	return true
}

func Route(action string, message Message) {
	defer func() {
		if err := recover(); err != nil {
			ErrorResponse(message, "invalid")
		}
	}()

	Routes[message.Action](message)
}

func Initialize() {
	Routes = make(map[string]func(Message))
}

func TestConnection() {
	go func() {
		for {
			time.Sleep(time.Second * 5)

			// Send ping
			send.Pipe(pipe.Message{
				Channel: pipe.BroadcastChannel([]int64{1, 3}),
				Event: pipe.Event{
					Sender: 0,
					Name:   "ping",
					Data: map[string]interface{}{
						"node": pipe.CurrentNode.ID,
					},
				},
			})
		}
	}()
}
