package handler

import (
	"chat-node/bridge"
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
