package message

import "chat-node/handler"

func SetupActions() {
	handler.Routes["conv_msg_create"] = createMessage
	handler.Routes["conv_msg_update"] = updateMessage

	// Typing status
	handler.Routes["conv_t_s"] = typingStatus
	handler.Routes["conv_t"] = typingStatus
}
