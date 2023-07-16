package service

import (
	"chat-node/database"
	"chat-node/database/fetching"
	"chat-node/util"

	"github.com/Fajurion/pipes"
	"github.com/Fajurion/pipesfiber"
)

type statusEntity struct {
	ID     string `json:"account"`
	Status string `json:"status"`
	Type   uint   `json:"type"`
}

// Setup the friends of the user (online)
func setup_fr(client *pipesfiber.Client, account *string, current *fetching.Session) bool {

	// Get the friends of the user
	res, err := util.PostRequest("/account/friends/online", map[string]interface{}{
		"node":    util.NODE_ID,
		"token":   util.NODE_TOKEN,
		"account": account,
	})

	if err != nil {
		return false
	}

	// Get status of friends
	var status []statusEntity
	if database.DBConn.Model(&fetching.Status{}).Select("id,status,type").Where("id IN ?", res["friends"]).Scan(&status).Error != nil {
		return false
	}

	// Get status of the user
	var userStatus fetching.Status
	if database.DBConn.Model(&fetching.Status{}).Where("id = ?", account).Take(&userStatus).Error != nil {
		return false
	}

	// Send the friends to the user
	client.SendEvent(pipes.Event{
		Name: "setup_st",
		Data: map[string]interface{}{
			"status": status,
			"own_status": map[string]interface{}{
				"status": userStatus.Status,
				"type":   userStatus.Type,
			},
		},
	})

	return true
}
