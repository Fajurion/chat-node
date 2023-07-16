package message

import "github.com/Fajurion/pipesfiber/wshandler"

func SetupActions() {
	wshandler.Routes["conv_msg_create"] = createMessage
	wshandler.Routes["conv_msg_update"] = updateMessage

	// Typing status
	wshandler.Routes["conv_t_s"] = typingStatus
	wshandler.Routes["conv_t"] = typingStatus
}
