package friends

import (
	"chat-node/util"

	"github.com/Fajurion/pipes"
	"github.com/Fajurion/pipes/send"
	"github.com/Fajurion/pipesfiber/wshandler"
	"github.com/gofiber/fiber/v2"
)

// Action: fr_rq_deny
func denyFriendRequest(message wshandler.Message) {

	if message.ValidateForm("id") {
		wshandler.ErrorResponse(message, "invalid")
		return
	}

	id := int64(message.Data["id"].(float64))

	res, err := util.PostRequest("/account/friends/request/deny", fiber.Map{
		"id":      util.NODE_ID,
		"token":   util.NODE_TOKEN,
		"session": message.Client.Session,
		"account": id,
	})

	if err != nil {
		return
	}

	success := res["success"].(bool)

	if !success {
		wshandler.ErrorResponse(message, res["error"].(string))
		return
	}

	if res["action"] == nil {
		wshandler.SuccessResponse(message)
		return
	}

	action := res["action"].(string)

	if action != "deny" {
		wshandler.ErrorResponse(message, "server.error")
		return
	}

	nodeRaw := res["node"].(map[string]interface{})
	nodeEntity := pipes.Node{
		ID:    util.Node64(int64(nodeRaw["id"].(float64))),
		WS:    nodeRaw["domain"].(string),
		Token: nodeRaw["token"].(string),
	}
	friend := res["friend"].(string)
	userData := message.Client.Data.(util.UserData)

	send.Socketless(nodeEntity, pipes.Message{
		Channel: pipes.BroadcastChannel([]string{friend}),
		Event: pipes.Event{
			Sender: message.Client.ID,
			Name:   "fr_rq:l",
			Data: map[string]interface{}{
				"status":   "denied",
				"username": userData.Username,
				"tag":      userData.Tag,
				"id":       message.Client.ID,
			},
		},
	})

	wshandler.StatusResponse(message, "denied")
}
