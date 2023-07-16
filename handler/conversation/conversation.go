package conversation

import (
	"chat-node/handler/conversation/call"
	"chat-node/handler/conversation/message"

	"github.com/Fajurion/pipesfiber/wshandler"
)

func SetupActions() {
	wshandler.Routes["conv_open"] = openConversation
	wshandler.Routes["conv_mem"] = getConversationMembers

	message.SetupActions()
	call.SetupActions()
}
