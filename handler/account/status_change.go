package account

import (
	"chat-node/database"
	"chat-node/database/fetching"
	"chat-node/util/localization"

	"github.com/Fajurion/pipes"
	"github.com/Fajurion/pipes/send"
	"github.com/Fajurion/pipesfiber/wshandler"
)

// Action: st_ch
func changeStatus(message wshandler.Message) {

	if !message.ValidateForm("status") {
		return
	}
	status := message.Data["status"].(string)

	// Save in database
	if err := database.DBConn.Model(&fetching.Status{}).Where("id = ?", message.Client.ID).Update("data", status).Error; err != nil {
		wshandler.ErrorResponse(message, localization.ErrorServer)
		return
	}

	// Send to all clients
	send.Client(message.Client.ID, pipes.Event{
		Name: "o:acc_st", // o: for own
		Data: map[string]interface{}{
			"st": status,
		},
	})

	wshandler.SuccessResponse(message)
}
