package conversation

import (
	"chat-node/handler"
	"chat-node/handler/conversation/message"
)

func SetupActions() {
	handler.Routes["conv_open"] = openConversation

	message.SetupActions()
}
