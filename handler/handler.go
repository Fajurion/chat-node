package handler

import (
	"chat-node/handler/account"

	"github.com/Fajurion/pipesfiber/wshandler"
)

func Create() {
	wshandler.Initialize()

	account.SetupActions()
}
