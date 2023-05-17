package service

import (
	"chat-node/bridge"
	"chat-node/database"
	"chat-node/database/fetching"
	"log"

	"github.com/Fajurion/pipes"
)

func setup_act(client *bridge.Client, current *fetching.Session, firstFetch *int64) bool {

	// Delete outdated stored actions
	database.DBConn.Where("created_at < ? AND account = ?", *firstFetch, client.ID).Delete(&fetching.Action{})

	// Get new actions
	var actionList []fetching.Action
	database.DBConn.Model(&fetching.Action{}).Where("created_at > ? AND account = ?", current.LastFetch, current.Account).Take(&actionList)

	log.Println("Actions: ", actionList)

	// Send the actions to the user
	client.SendEvent(pipes.Event{
		Name: "setup_act",
		Data: map[string]interface{}{
			"actions": actionList,
		},
	})

	return true
}
