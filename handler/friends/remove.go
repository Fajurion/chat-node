package friends

import (
	"chat-node/handler"
)

// Action: fr_rem
func removeFriend(message handler.Message) {
	handler.ErrorResponse(message, "not.implemented")

	/*
		// Check if request is valid
		if message.ValidateForm("id") {
			handler.ErrorResponse(message, "invalid")
			return
		}

		// Get the account
		id := int64(message.Data["id"].(float64))

		// Send request to the server
		res, err := util.PostRequest("/account/friends/remove", fiber.Map{
			"node":    util.NODE_ID,
			"token":   util.NODE_TOKEN,
			"account": message.Client.ID,
			"friend":  id,
		})

		if err != nil {
			handler.ErrorResponse(message, "server.error")
			return
		}

		if !res["success"].(bool) {
			handler.ErrorResponse(message, res["error"].(string))
			return
		}

		//* Create stored action
		database.DBConn.Create(&fetching.Action{
			ID:      util.GenerateToken(32),
			Account: message.Client.ID,
			Action:  "fr_rem",
			Target:  fmt.Sprintf("%d", id),
		})

		database.DBConn.Create(&fetching.Action{
			ID:      util.GenerateToken(32),
			Account: id,
			Action:  "fr_rem",
			Target:  fmt.Sprintf("%d", message.Client.ID),
		})

		handler.SuccessResponse(message) */
}
