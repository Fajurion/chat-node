package call

import (
	"github.com/Fajurion/pipesfiber/wshandler"
)

func SetupActions() {
	wshandler.Routes["c_s"] = start
	wshandler.Routes["c_j"] = join
	wshandler.Routes["c_c"] = status
}
