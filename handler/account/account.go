package account

import "chat-node/handler"

func SetupActions() {
	handler.Routes["acc_st"] = changeStatus
	handler.Routes["acc_on"] = setOnline
}
