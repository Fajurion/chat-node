package account

import "github.com/Fajurion/pipesfiber/wshandler"

func SetupActions() {
	wshandler.Routes["acc_st"] = changeStatus
	wshandler.Routes["acc_on"] = setOnline
}
