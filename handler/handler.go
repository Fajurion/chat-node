package handler

import (
	"chat-node/bridge"
	"chat-node/pipe"
	"chat-node/pipe/send"
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

	// Check if the action exists
	if Routes[message.Action] == nil {
		return false
	}

	go Routes[message.Action](message)

	return true
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
				Channel: pipe.BroadcastChannel(1, []int64{2}),
				Event: pipe.Event{
					Name: "ping",
					Data: map[string]interface{}{},
				},
			})
		}
	}()
}
