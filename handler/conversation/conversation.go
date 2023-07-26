package conversation

import (
	"chat-node/handler/conversation/message"

	"github.com/Fajurion/pipesfiber/wshandler"
)

func SetupActions() {
	wshandler.Routes["conv_open"] = openConversation

	message.SetupActions()

	// TODO: Implement new call system
	//call.SetupActions()
}
