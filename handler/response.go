package handler

import (
	"chat-node/util"
	"runtime/debug"

	"github.com/Fajurion/pipes"
	"github.com/Fajurion/pipes/send"
)

func NormalResponse(message Message, data map[string]interface{}) {
	Response(message.Client.ID, message.Action, data)
}

func Response(client string, action string, data map[string]interface{}) {
	send.Client(client, pipes.Event{
		Sender: "0",
		Name:   action,
		Data:   data,
	})
}

func SuccessResponse(message Message) {
	Response(message.Client.ID, message.Action, map[string]interface{}{
		"success": true,
		"message": message,
	})
}

func StatusResponse(message Message, status string) {
	Response(message.Client.ID, message.Action, map[string]interface{}{
		"success": true,
		"message": status,
	})
}

func ErrorResponse(message Message, err string) {

	if util.Testing {
		debug.PrintStack()
	}

	Response(message.Client.ID, message.Action, map[string]interface{}{
		"success": false,
		"message": err,
	})
}

func (message *Message) ValidateForm(fields ...string) bool {

	for _, field := range fields {
		if message.Data[field] == nil {
			return true
		}
	}

	return false
}
