package friends

import (
	"chat-node/handler"
	"chat-node/util"

	"github.com/Fajurion/pipes"
	"github.com/Fajurion/pipes/send"
	"github.com/gofiber/fiber/v2"
)

// Action: fr_rq_deny
func denyFriendRequest(message handler.Message) {

	if message.ValidateForm("id") {
		handler.ErrorResponse(message, "invalid")
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
		handler.ErrorResponse(message, res["error"].(string))
		return
	}

	if res["action"] == nil {
		handler.SuccessResponse(message)
		return
	}

	action := res["action"].(string)

	if action != "deny" {
		handler.ErrorResponse(message, "server.error")
		return
	}

	nodeRaw := res["node"].(map[string]interface{})
	nodeEntity := pipes.Node{
		ID:    util.User64(int64(nodeRaw["id"].(float64))),
		WS:    nodeRaw["domain"].(string),
		Token: nodeRaw["token"].(string),
	}
	friend := int64(res["friend"].(float64))

	send.Socketless(nodeEntity, pipes.Message{
		Channel: pipes.BroadcastChannel([]string{util.User64(friend)}),
		Event: pipes.Event{
			Sender: util.User64(message.Client.ID),
			Name:   "fr_rq:l",
			Data: map[string]interface{}{
				"status":   "denied",
				"username": message.Client.Username,
				"tag":      message.Client.Tag,
				"id":       message.Client.ID,
			},
		},
	})

	handler.StatusResponse(message, "denied")
}
