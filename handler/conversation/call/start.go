package call

import (
	"chat-node/calls"
	"context"
	"fmt"

	"github.com/Fajurion/pipesfiber/wshandler"
	"github.com/livekit/protocol/livekit"
)

func start(message wshandler.Message) {

	if message.ValidateForm("id") {
		wshandler.ErrorResponse(message, "invalid")

		println("invalid form")
		return
	}

	id := int64(message.Data["id"].(float64)) // ID of the conversation
	roomName := fmt.Sprintf("c_%d", id)       // Livekit room name

	// Check if user is member of the conversation
	/*
		if database.DBConn.Where("conversation = ? AND account = ?", id, message.Client.ID).Find(&conversations.Member{}).Error != nil {
			wshandler.ErrorResponse(message, "invalid")
			return
		} */

	// Check if there is already a call (livekit room)
	res, _ := calls.RoomClient.ListRooms(context.Background(), &livekit.ListRoomsRequest{
		Names: []string{roomName},
	})

	if len(res.Rooms) > 0 {

		println("connecting to call")

		// Connect to call
		tk, err := calls.GetJoinToken(roomName, message.Client.ID)

		if err != nil {
			wshandler.ErrorResponse(message, "server.error")
			return
		}

		wshandler.NormalResponse(message, map[string]interface{}{
			"success": true,
			"call":    true,
			"id":      id,
			"token":   tk,
		})

		return
	}

	// Create call invite
	token, err := calls.GenerateCallToken(roomName, message.Client.ID)

	if err != nil {
		wshandler.ErrorResponse(message, "server.error")
		return
	}

	wshandler.NormalResponse(message, map[string]interface{}{
		"success": true,
		"call":    false,
		"id":      id,
		"token":   token,
	})
}
