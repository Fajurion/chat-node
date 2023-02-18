package handler

import (
	"chat-node/pipe"
	"chat-node/pipe/send"
)

func Response(client int64, action string, data map[string]interface{}) {
	send.Pipe(pipe.Message{
		Channel: pipe.ClientChannel(client),
		Event: pipe.Event{
			Name: action,
			Data: data,
		},
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
	Response(message.Client.ID, message.Action, map[string]interface{}{
		"success": false,
		"message": err,
	})
}