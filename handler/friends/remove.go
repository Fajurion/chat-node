package friends

import (
	"chat-node/database"
	"chat-node/database/fetching"
	"chat-node/util"

	"github.com/Fajurion/pipesfiber/wshandler"
	"github.com/gofiber/fiber/v2"
)

// Action: fr_rem
func removeFriend(message wshandler.Message) {

	// Check if request is valid
	if message.ValidateForm("id") {
		wshandler.ErrorResponse(message, "invalid")
		return
	}

	// Get the account
	id := message.Data["id"].(string)

	// Send request to the server
	res, err := util.PostRequest("/account/friends/remove", fiber.Map{
		"node":    util.NODE_ID,
		"token":   util.NODE_TOKEN,
		"account": message.Client.ID,
		"friend":  id,
	})

	if err != nil {
		wshandler.ErrorResponse(message, "server.error")
		return
	}

	if !res["success"].(bool) {
		wshandler.ErrorResponse(message, res["error"].(string))
		return
	}

	//* Create stored action
	database.DBConn.Create(&fetching.Action{
		ID:      util.GenerateToken(32),
		Account: message.Client.ID,
		Action:  "fr_rem",
		Target:  id,
	})

	database.DBConn.Create(&fetching.Action{
		ID:      util.GenerateToken(32),
		Account: id,
		Action:  "fr_rem",
		Target:  message.Client.ID,
	})

	wshandler.SuccessResponse(message)
}
