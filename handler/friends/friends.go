package friends

import (
	"chat-node/handler"
)

func SetupActions() {
	handler.Routes["friend_request"] = friendRequest
	handler.Routes["friend_request_deny"] = denyFriendRequest
}
