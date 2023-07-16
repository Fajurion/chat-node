package handler

import (
	"chat-node/handler/account"
	"chat-node/handler/conversation"
	"chat-node/handler/friends"

	"github.com/Fajurion/pipesfiber/wshandler"
)

func Create() {
	wshandler.Initialize()

	friends.SetupActions()
	conversation.SetupActions()
	account.SetupActions()
}
