package conversation

import (
	"chat-node/handler"
	"chat-node/handler/conversation/call"
	"chat-node/handler/conversation/message"
)

func SetupActions() {
	handler.Routes["conv_open"] = openConversation
	handler.Routes["conv_mem"] = getConversationMembers

	message.SetupActions()
	call.SetupActions()
}
