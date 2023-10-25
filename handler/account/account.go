package account

import "github.com/Fajurion/pipesfiber/wshandler"

func SetupActions() {
	wshandler.Routes["st_ch"] = changeStatus
	wshandler.Routes["st_send"] = sendStatus
	wshandler.Routes["st_res"] = respondToStatus
}
