package call

import (
	"chat-node/handler"
)

func SetupActions() {
	handler.Routes["c_s"] = start
	handler.Routes["c_j"] = join
}
