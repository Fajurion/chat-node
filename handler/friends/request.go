package friends

import (
	"chat-node/handler"
	"chat-node/pipe"
	"chat-node/pipe/send"
	"chat-node/util"

	"github.com/gofiber/fiber/v2"
)

// Action: friend_request
func friendRequest(message handler.Message) {
	if message.ValidateForm("username", "tag") {
		handler.ErrorResponse(message, "invalid")
		return
	}

	username := message.Data["username"].(string)
	tag := message.Data["tag"].(string)

	res, err := util.PostRequest("/account/friends/request/create", fiber.Map{
		"id":       util.NODE_ID,
		"token":    util.NODE_TOKEN,
		"session":  message.Client.Session,
		"username": username,
		"tag":      tag,
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
	nodeRaw := res["node"].(map[string]interface{})
	nodeEntity := pipe.Node{
		ID:     int64(nodeRaw["id"].(float64)),
		Domain: nodeRaw["domain"].(string),
		Token:  nodeRaw["token"].(string),
		App:    uint(nodeRaw["id"].(float64)),
	}
	friend := int64(res["friend"].(float64))

	switch action {
	case "accept":
		handler.StatusResponse(message, "accepted")

		send.Socketless(nodeEntity, pipe.Message{
			Channel: pipe.BroadcastChannel([]int64{friend}),
			Event: pipe.Event{
				Sender: message.Client.ID,
				Name:   "friend_request",
				Data: map[string]interface{}{
					"status":   "accepted",
					"username": message.Client.Username,
					"tag":      message.Client.Tag,
					"id":       message.Client.ID,
				},
			},
		})

	case "send":
		handler.StatusResponse(message, "sent")

		send.Socketless(nodeEntity, pipe.Message{
			Channel: pipe.BroadcastChannel([]int64{friend}),
			Event: pipe.Event{
				Sender: message.Client.ID,
				Name:   "friend_request",
				Data: map[string]interface{}{
					"status":   "sent",
					"username": message.Client.Username,
					"tag":      message.Client.Tag,
					"id":       message.Client.ID,
				},
			},
		})

	}
}
