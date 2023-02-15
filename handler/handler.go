package handler

import (
	"chat-node/bridge"
	"chat-node/pipe"
	"chat-node/pipe/send"
	"time"
)

type Message struct {
	Client bridge.Client          `json:"client"`
	Data   map[string]interface{} `json:"data"`
}

// Routes is a map of all the routes
var Routes map[string]func(Message) error

func Handle(action string, message Message) bool {

	// Check if the action exists
	if Routes[action] == nil {
		return false
	}

	go Routes[action](message)

	return true
}

func Initialize() {
	Routes = make(map[string]func(Message) error)
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
