package processors

import (
	"chat-node/pipe"
	"log"

	"github.com/bytedance/sonic"
)

var Processors map[string]func(*pipe.Message, int64) pipe.Event = make(map[string]func(*pipe.Message, int64) pipe.Event)

func ProcessMarshal(message *pipe.Message, target int64) []byte {
	event := ProcessEvent(message, target)

	// Marshal the event
	msg, err := sonic.Marshal(event)
	if err != nil {
		return nil
	}

	return msg
}

func ProcessEvent(message *pipe.Message, target int64) pipe.Event {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("Error processing message: %s \n", err)
		}
	}()

	if Processors[message.Event.Name] != nil {
		return Processors[message.Event.Name](message, target)
	}

	return message.Event
}
