package service

import (
	"chat-node/bridge"
	"chat-node/database"
	"chat-node/database/fetching"
	"chat-node/pipe"
	"chat-node/util"
)

type statusEntity struct {
	ID     int64  `json:"account"`
	Status string `json:"status"`
}

// Setup the friends of the user (online)
func setup_fr(client *bridge.Client, account *int64, current *fetching.Session) bool {

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
	if database.DBConn.Model(&fetching.Status{}).Select("id,status").Where("id IN ?", res["friends"]).Scan(&status).Error != nil {
		return false
	}

	// Send the friends to the user
	client.SendEvent(pipe.Event{
		Name: "setup_st",
		Data: map[string]interface{}{
			"status": status,
		},
	})

	return true
}