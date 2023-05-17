package account

import (
	"chat-node/bridge"
	"chat-node/database"
	"chat-node/database/fetching"
	"chat-node/handler"
	"chat-node/util"
	"fmt"

	"github.com/Fajurion/pipes"
	"github.com/Fajurion/pipes/send"
)

// Action: acc_st
func changeStatus(message handler.Message) {

	if message.ValidateForm("status") {
		handler.ErrorResponse(message, "invalid")
		return
	}

	valid, err := UpdateStatus(message.Client, fetching.StatusOnline, message.Data["status"].(string), true)

	if !valid {
		handler.ErrorResponse(message, err)
		return
	}

	handler.SuccessResponse(message)
}

// Action: acc_on
func setOnline(message handler.Message) {

	valid, err := UpdateStatus(message.Client, fetching.StatusOnline, "", false)

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
		database.DBConn.Model(&fetching.Status{}).Where("id = ?", client.ID).Update("status", status)
	} else {
		database.DBConn.Model(&fetching.Status{}).Select("status").Where("id = ?", client.ID).Scan(&status)
	}

	// Transform array
	var friends []string
	for _, friend := range res["friends"].([]interface{}) {
		friends = append(friends, fmt.Sprintf("%d", uint(friend.(float64))))
	}

	// Send the event to the friends
	send.Pipe(send.ProtocolWS, pipes.Message{
		Event: pipes.Event{
			Name:   "fr_st",
			Sender: util.User64(client.ID),
			Data: map[string]interface{}{
				"st": status,
			},
		},
		Channel: pipes.BroadcastChannel(friends),
	})

	// Send the response
	return true, ""
}
