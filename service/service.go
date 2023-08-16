package service

import (
	"chat-node/database"
	"chat-node/database/fetching"
	"chat-node/util"

	integration "fajurion.com/node-integration"
	"github.com/Fajurion/pipes"
	"github.com/Fajurion/pipesfiber"
)

func User(client *pipesfiber.Client) bool {
	account := client.ID

	// Check if the account is already in the database
	var status fetching.Status
	if database.DBConn.Where(&fetching.Status{ID: account}).Take(&status).Error != nil {

		// Create a new status
		if database.DBConn.Create(&fetching.Status{
			ID:   account,
			Data: "-", // Status is disabled
			Node: integration.NODE_ID,
		}).Error != nil {
			return false
		}
	} else {

		// Update the status
		database.DBConn.Model(&fetching.Status{}).Where("id = ?", account).Update("node", util.NodeTo64(pipes.CurrentNode.ID))
	}

	// Send current status
	client.SendEvent(pipes.Event{
		Name: "setup_st", // :n = new
		Data: map[string]interface{}{
			"data": status.Data,
			"node": status.Node,
		},
	})

	// Send the setup complete event
	client.SendEvent(pipes.Event{
		Name: "setup_fin",
		Data: map[string]interface{}{},
	})

	return true
}
