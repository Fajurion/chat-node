package conversation

import (
	"chat-node/handler/conversation/space"
	"log"
	"time"

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
	tokenID string
	token   string
	date    int64 // Unix timestamp of last fetch
}

var newTaskChan = make(chan TokenTask)

func setupMessageQueue() {
	for i := 0; i < 1; i++ {
		go func() {
			for {

				// Wait for a new task
				<-newTaskChan
				log.Println("doing task")
				time.Sleep(500 * time.Millisecond)

				// TODO: Fetch all messages from the database
			}
		}()
	}
}

func AddConversationToken(task TokenTask) {
	newTaskChan <- task
}
