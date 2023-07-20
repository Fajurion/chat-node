package service

import (
	"chat-node/database"
	"chat-node/database/fetching"
	"chat-node/util"
	"log"
	"time"

	integration "fajurion.com/node-integration"
	"github.com/Fajurion/pipes"
	"github.com/Fajurion/pipesfiber"
)

func User(client *pipesfiber.Client) bool {
	session := client.Session
	account := client.ID
	data := client.Data.(util.UserData)

	client.SendEvent(pipes.Event{
		Name: "setup_wel",
		Data: map[string]interface{}{
			"name": data.Username,
			"tag":  data.Tag,
		},
	})

	// Check if the account is already in the database
	var status fetching.Status
	if database.DBConn.Where(&fetching.Status{ID: account}).Take(&status).Error != nil {

		// Create a new status
		if database.DBConn.Create(&fetching.Status{
			ID:     account,
			Type:   fetching.StatusOnline,
			Status: "-",
			Node:   util.NodeTo64(pipes.CurrentNode.ID),
		}).Error != nil {
			return false
		}
	} else {

		// Update the status
		database.DBConn.Model(&fetching.Status{}).Where(&fetching.Status{ID: account}).Update("node", util.NodeTo64(pipes.CurrentNode.ID))
	}

	// Check if this is a new device
	var current fetching.Session
	if database.DBConn.Where(&fetching.Session{ID: session}).Take(&current).Error != nil {

		// TODO: New device (sync with old device)

		// Save the session
		current = fetching.Session{
			ID:        session,
			Account:   account,
			Node:      util.NodeTo64(pipes.CurrentNode.ID),
			LastFetch: 0,
		}

		if database.DBConn.Create(&current).Error != nil {
			return false
		}

		client.SendEvent(pipes.Event{
			Name: "setup_device",
			Data: map[string]interface{}{
				"device": client.Session,
			},
		})
	}

	// Get the earliest fetch time
	var firstFetch int64
	database.DBConn.Raw("SELECT MIN(last_fetch) FROM sessions WHERE account = ?", account).Scan(&firstFetch)

	// Get new conversations
	if !setup_conv(client, &account, &current) {
		return false
	}

	if integration.Testing {
		log.Println("Fetch:", current.LastFetch)
	}

	// Get new actions
	if !setup_act(client, &current, &firstFetch) {
		return false
	}

	if integration.Testing {
		log.Println("setup_act finished for", client.ID)
	}

	// Get new messages
	if !setup_mes(client, &current, &account) {
		return false
	}

	if integration.Testing {
		log.Println("setup_mes finished for", client.ID)
	}

	// Get online friends
	if err := setup_fr(client, &account, &current); err != nil {
		log.Println(err)
		return false
	}

	if integration.Testing {
		log.Println("setup_fr finished for", client.ID)
	}

	// Save the session
	current.LastFetch = time.Now().UnixMilli()
	current.Node = util.NodeTo64(pipes.CurrentNode.ID)
	if database.DBConn.Save(&current).Error != nil {
		pipesfiber.Remove(client.ID, client.Session)
		return false
	}

	// Send the setup complete event
	client.SendEvent(pipes.Event{
		Name: "setup_fin",
		Data: map[string]interface{}{},
	})

	return true
}
