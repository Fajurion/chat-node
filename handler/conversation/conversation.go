package conversation

import (
	"chat-node/database"
	"chat-node/database/conversations"
	"chat-node/handler/conversation/space"
	message_routes "chat-node/routes/conversations/message"

	"github.com/Fajurion/pipes/send"
	"github.com/Fajurion/pipesfiber/wshandler"
)

func SetupActions() {
	space.SetupActions()

	wshandler.Routes["conv_sub"] = subscribe

	// Setup messages queue
	setupMessageQueue()
}

const messageProcessorAmount = 3

type TokenTask struct {
	Adapter      string
	Conversation string
	Date         int64 // Unix timestamp of last fetch
}

var newTaskChan = make(chan TokenTask)

func setupMessageQueue() {
	for i := 0; i < messageProcessorAmount; i++ {
		go func() {
			for {

				// Wait for a new task
				task := <-newTaskChan

				// Get all messages
				var messages []conversations.Message
				if database.DBConn.Where("conversation = ? AND creation > ?", task.Conversation, task.Date).Find(&messages).Error != nil {
					continue
				}

				// Send messages to the adapter
				for _, message := range messages {
					send.Client(task.Adapter, message_routes.MessageEvent(message))
				}
			}
		}()
	}
}

func AddConversationToken(task TokenTask) {
	newTaskChan <- task
}
