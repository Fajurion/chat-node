package friends

import "github.com/Fajurion/pipesfiber/wshandler"

func SetupActions() {
	wshandler.Routes["fr_rq"] = friendRequest
	wshandler.Routes["fr_rq_deny"] = denyFriendRequest

	wshandler.Routes["fr_rem"] = removeFriend
}
