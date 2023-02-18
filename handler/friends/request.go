package friends

import (
	"chat-node/handler"
	"chat-node/pipe"
	"chat-node/pipe/send"
	"chat-node/util"

	"github.com/gofiber/fiber/v2"
)

func friendRequest(message handler.Message) {

	username := message.Data["username"].(string)
	tag := message.Data["tag"].(string)

	res, err := util.PostRequest("/account/friends/request/create", fiber.Map{
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
		handler.ErrorResponse(message, res["message"].(string))
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

		send.Pipe(pipe.Message{
			Channel: pipe.SocketlessChannel(message.Client.ID, nodeEntity, []int64{friend}),
			Event: pipe.Event{
				Name: "friend_request",
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

		send.Pipe(pipe.Message{
			Channel: pipe.SocketlessChannel(message.Client.ID, nodeEntity, []int64{friend}),
			Event: pipe.Event{
				Name: "friend_request",
				Data: map[string]interface{}{
					"status":   "sent",
					"username": message.Client.Username,
					"tag":      message.Client.Tag,
					"id":       message.Client.ID,
				},
			},
		})

		// TODO: Send notification to user
	}
}
