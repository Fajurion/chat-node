package setup

import (
	"chat-node/handler"
	"chat-node/handler/account"
	"chat-node/handler/friends"
	"chat-node/handler/project"
)

func Initialize() {
	handler.Initialize()

	friends.SetupActions()
	project.SetupActions()
	account.SetupActions()
}
