package friends

import (
	"chat-node/handler"
	"chat-node/pipe"
	"chat-node/pipe/send"
	"chat-node/util"

	"github.com/Fajurion/pipes"
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
		ID:     int64(nodeRaw["id"].(float64)),
		Domain: nodeRaw["domain"].(string),
		Token:  nodeRaw["token"].(string),
		App:    uint(nodeRaw["id"].(float64)),
	}
	friend := int64(res["friend"].(float64))

	send.Socketless(nodeEntity, pipe.Message{
		Channel: pipe.BroadcastChannel([]int64{friend}),
		Event: pipe.Event{
			Sender: message.Client.ID,
			Name:   "friend_request",
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
