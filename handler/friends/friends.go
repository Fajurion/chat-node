package friends

import (
	"chat-node/handler"
)

func SetupActions() {
	handler.Routes["friends"] = friendRequest
}
