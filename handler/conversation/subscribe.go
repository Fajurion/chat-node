package conversation

import (
	"chat-node/database"
	"chat-node/database/fetching"
	"chat-node/util"

	"github.com/Fajurion/pipes"
	"github.com/Fajurion/pipes/adapter"
	"github.com/Fajurion/pipesfiber/wshandler"
	"gorm.io/gorm/clause"
)

// Action: conv_sub
func subscribe(message wshandler.Message) {

	if message.ValidateForm("tokens") {
		wshandler.ErrorResponse(message, "invalid")
		return
	}

	tokens := message.Data["tokens"].([]interface{})

	// TODO: Validate tokens

	for _, token := range tokens {

		// Register adapter for the subscription
		adapter.AdaptWS(adapter.Adapter{
			ID: "s-" + token.(string),
			Receive: func(ctx *adapter.Context) error {
				return message.Client.SendEvent(*ctx.Event)
			},
		})

		// Register the subscription in the database
		database.DBConn.Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Create(&fetching.Subscription{
			ID:   token.(string),
			Node: util.NodeTo64(pipes.CurrentNode.ID),
		})
	}

	wshandler.SuccessResponse(message)
}
