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
	session := client.Session
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
		database.DBConn.Model(&fetching.Status{}).Where(&fetching.Status{ID: account}).Update("node", util.NodeTo64(pipes.CurrentNode.ID))
	}

	// Send current status
	client.SendEvent(pipes.Event{
		Name: "setup_st", // :n = new
		Data: map[string]interface{}{
			"data": status.Data,
			"node": status.Node,
		},
	})

	// Check if the account already has a mailbox
	var current fetching.Mailbox
	if database.DBConn.Where(&fetching.Mailbox{ID: account}).Take(&current).Error != nil {

		// TODO: New device (sync with old device)

		// Save the session
		current = fetching.Mailbox{
			ID:    session,
			Token: util.GenerateToken(util.MailboxTokenLength),
		}

		if database.DBConn.Create(&current).Error != nil {
			return false
		}

		client.SendEvent(pipes.Event{
			Name: "setup_mail:n", // :n = new
			Data: map[string]interface{}{
				"new":   true,
				"token": current.Token,
			},
		})
	} else {
		client.SendEvent(pipes.Event{
			Name: "setup_mail",
			Data: map[string]interface{}{
				"token": current.Token,
			},
		})
	}

	// Send the setup complete event
	client.SendEvent(pipes.Event{
		Name: "setup_fin",
		Data: map[string]interface{}{},
	})

	return true
}
