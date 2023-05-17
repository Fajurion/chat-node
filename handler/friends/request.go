package friends

import (
	"chat-node/database"
	"chat-node/database/fetching"
	"chat-node/handler"
	"chat-node/util"

	"github.com/Fajurion/pipes"
	"github.com/Fajurion/pipes/send"
	"github.com/gofiber/fiber/v2"
)

// Action: fr_rq
func friendRequest(message handler.Message) {

	if message.ValidateForm("username", "tag", "signature") {
		handler.ErrorResponse(message, "invalid")
		return
	}

	username := message.Data["username"].(string)
	tag := message.Data["tag"].(string)

	res, err := util.PostRequest("/account/friends/request/create", fiber.Map{
		"id":        util.NODE_ID,
		"token":     util.NODE_TOKEN,
		"session":   message.Client.Session,
		"username":  username,
		"tag":       tag,
		"signature": message.Data["signature"].(string),
	})

	if err != nil {
		handler.ErrorResponse(message, "invalid")
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
	nodeEntity := pipes.Node{
		ID:    util.User64(int64(nodeRaw["id"].(float64))),
		WS:    nodeRaw["domain"].(string),
		Token: nodeRaw["token"].(string),
	}
	friend := int64(res["friend"].(float64))
	signature := res["signature"].(string)
	key := res["key"].(string)

	switch action {
	case "accept":

		// Delete stored actions for friend removal
		database.DBConn.
			Where("action = ? AND account IN ? AND target IN ?", "fr_rem", []int64{message.Client.ID, friend}, []int64{message.Client.ID, friend}).
			Delete(&fetching.Action{})

		handler.NormalResponse(message, map[string]interface{}{
			"success": true,
			"message": "accepted",
			"name":    username,
			"tag":     tag,
			"id":      friend,
		})

		send.Socketless(nodeEntity, pipes.Message{
			Channel: pipes.BroadcastChannel([]string{util.User64(friend)}),
			Event: pipes.Event{
				Sender: util.User64(message.Client.ID),
				Name:   "fr_rq:l",
				Data: map[string]interface{}{
					"status":    "accepted",
					"name":      message.Client.Username,
					"tag":       message.Client.Tag,
					"id":        message.Client.ID,
					"signature": signature,
					"key":       key,
				},
			},
		})

	case "send":
		handler.NormalResponse(message, map[string]interface{}{
			"success": true,
			"message": "sent",
			"name":    username,
			"tag":     tag,
			"id":      friend,
		})

		send.Socketless(nodeEntity, pipes.Message{
			Channel: pipes.BroadcastChannel([]string{util.User64(friend)}),
			Event: pipes.Event{
				Sender: util.User64(message.Client.ID),
				Name:   "fr_rq:l",
				Data: map[string]interface{}{
					"status":    "sent",
					"name":      message.Client.Username,
					"tag":       message.Client.Tag,
					"id":        message.Client.ID,
					"signature": signature,
					"key":       key,
				},
			},
		})
	}
}
