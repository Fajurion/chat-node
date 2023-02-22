package handler

import (
	"chat-node/pipe"
	"chat-node/pipe/send"
)

func Response(client int64, action string, data map[string]interface{}) {
	send.Client(client, pipe.Event{
		Sender: 0,
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
