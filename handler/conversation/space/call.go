package space

import (
	"github.com/Fajurion/pipesfiber/wshandler"
)

func SetupActions() {
	wshandler.Routes["spc_start"] = start
	wshandler.Routes["spc_join"] = joinCall
	wshandler.Routes["spc_leave"] = leaveCall
}
