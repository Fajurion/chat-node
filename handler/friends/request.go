package friends

import (
	"chat-node/database"
	"chat-node/database/fetching"
	"chat-node/util"

	"github.com/Fajurion/pipes"
	"github.com/Fajurion/pipes/send"
	"github.com/Fajurion/pipesfiber/wshandler"
	"github.com/gofiber/fiber/v2"
)

// Action: fr_rq
func friendRequest(message wshandler.Message) {

	if message.ValidateForm("username", "tag", "signature") {
		wshandler.ErrorResponse(message, "invalid")
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
		wshandler.ErrorResponse(message, "invalid")
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
	nodeRaw := res["node"].(map[string]interface{})
	nodeEntity := pipes.Node{
		ID:    util.Node64(int64(nodeRaw["id"].(float64))),
		WS:    nodeRaw["domain"].(string),
		Token: nodeRaw["token"].(string),
	}
	friend := res["friend"].(string)
	signature := res["signature"].(string)
	key := res["key"].(string)
	userData := message.Client.Data.(util.UserData)

	switch action {
	case "accept":

		// Delete stored actions for friend removal
		database.DBConn.
			Where("action = ? AND account IN ? AND target IN ?", "fr_rem", []string{message.Client.ID, friend}, []string{message.Client.ID, friend}).
			Delete(&fetching.Action{})

		wshandler.NormalResponse(message, map[string]interface{}{
			"success": true,
			"message": "accepted",
			"name":    username,
			"tag":     tag,
			"id":      friend,
		})

		send.Socketless(nodeEntity, pipes.Message{
			Channel: pipes.BroadcastChannel([]string{friend}),
			Event: pipes.Event{
				Sender: message.Client.ID,
				Name:   "fr_rq:l",
				Data: map[string]interface{}{
					"status":    "accepted",
					"name":      userData.Username,
					"tag":       userData.Tag,
					"id":        message.Client.ID,
					"signature": signature,
					"key":       key,
				},
			},
		})

	case "send":
		wshandler.NormalResponse(message, map[string]interface{}{
			"success": true,
			"message": "sent",
			"name":    username,
			"tag":     tag,
			"id":      friend,
		})

		send.Socketless(nodeEntity, pipes.Message{
			Channel: pipes.BroadcastChannel([]string{friend}),
			Event: pipes.Event{
				Sender: message.Client.ID,
				Name:   "fr_rq:l",
				Data: map[string]interface{}{
					"status":    "sent",
					"name":      userData.Username,
					"tag":       userData.Tag,
					"id":        message.Client.ID,
					"signature": signature,
					"key":       key,
				},
			},
		})
	}
}
