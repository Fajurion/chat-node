package friends

import "chat-node/handler"

func SetupActions() {
	handler.Routes["friends"] = Handle
}

func Handle(message handler.Message) error {
	return nil
}
