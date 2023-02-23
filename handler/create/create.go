package create

import (
	"chat-node/handler"
	"chat-node/handler/account"
	"chat-node/handler/conversation"
	"chat-node/handler/friends"
)

func Create() {
	handler.Initialize()

	friends.SetupActions()
	conversation.SetupActions()
	account.SetupActions()
}
