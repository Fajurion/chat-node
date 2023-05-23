package account

import (
	"chat-node/bridge"
	"chat-node/database"
	"chat-node/database/fetching"
	"chat-node/handler"
	"chat-node/util"

	"github.com/Fajurion/pipes"
	"github.com/Fajurion/pipes/send"
)

// Action: acc_st
func changeStatus(message handler.Message) {

	if message.ValidateForm("type", "status") {
		handler.ErrorResponse(message, "invalid")
		return
	}

	sType := uint(message.Data["type"].(float64))
	if sType < fetching.StatusOnline || sType > fetching.StatusDoNotDisturb {
		handler.ErrorResponse(message, "invalid")
		return
	}

	valid, err := UpdateStatus(message.Client, sType, message.Data["status"].(string), true)

	if !valid {
		handler.ErrorResponse(message, err)
		return
	}

	handler.SuccessResponse(message)
}

// Action: acc_on
func setOnline(message handler.Message) {

	var sType uint = fetching.StatusOnline
	database.DBConn.Model(&fetching.Status{}).Select("type").Where("id = ?", message.Client.ID).Scan(&sType)
	valid, err := UpdateStatus(message.Client, sType, "", false)

	if !valid {
		handler.ErrorResponse(message, err)
		return
	}

	handler.SuccessResponse(message)
}

func UpdateStatus(client *bridge.Client, sType uint, status string, set bool) (bool, string) {

	// Send event through pipe
	res, err := util.PostRequest("/account/friends/online", map[string]interface{}{
		"node":    util.NODE_ID,
		"token":   util.NODE_TOKEN,
		"account": client.ID,
	})

	if err != nil {
		return false, "server.error"
	}

	// Update the status of the user
	if set {
		database.DBConn.Model(&fetching.Status{}).Where("id = ?", client.ID).Updates(map[string]interface{}{
			"type":   sType,
			"status": status,
		})
	} else {
		database.DBConn.Model(&fetching.Status{}).Select("status").Where("id = ?", client.ID).Scan(&status)
	}

	if res["friends"] == nil {
		return true, ""
	}

	// Transform array
	var friends []string
	for _, friend := range res["friends"].([]interface{}) {
		friends = append(friends, friend.(string))
	}

	if len(friends) == 0 {
		return true, ""
	}

	// Send the event to the friends
	send.Pipe(send.ProtocolWS, pipes.Message{
		Event: pipes.Event{
			Name:   "fr_st",
			Sender: client.ID,
			Data: map[string]interface{}{
				"t":  sType,
				"st": status,
			},
		},
		NoSelf:  true,
		Channel: pipes.BroadcastChannel(friends),
	})

	// Send the response
	return true, ""
}
