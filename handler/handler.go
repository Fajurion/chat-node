package handler

import (
	"chat-node/handler/account"
	"chat-node/handler/conversation"

	"github.com/Fajurion/pipesfiber/wshandler"
)

func Create() {
	wshandler.Initialize()

	conversation.SetupActions()
	account.SetupActions()
}
