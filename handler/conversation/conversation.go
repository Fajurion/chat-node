package conversation

import (
	"chat-node/handler/conversation/space"

	"github.com/Fajurion/pipesfiber/wshandler"
)

func SetupActions() {
	space.SetupActions()

	wshandler.Routes["conv_sub"] = subscribe
}
