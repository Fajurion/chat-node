package message

import "chat-node/handler"

func SetupActions() {
	handler.Routes["conv_msg_create"] = createMessage
	handler.Routes["conv_msg_update"] = updateMessage
}
