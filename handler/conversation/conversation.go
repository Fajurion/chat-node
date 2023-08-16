package conversation

import "github.com/Fajurion/pipesfiber/wshandler"

func SetupActions() {
	wshandler.Routes["conv_sub"] = subscribe
}
