package friends

import (
	"chat-node/handler"
)

func SetupActions() {
	handler.Routes["fr_rq"] = friendRequest
	handler.Routes["fr_rq_deny"] = denyFriendRequest
}
