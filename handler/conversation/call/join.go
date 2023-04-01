package call

import (
	"chat-node/calls"
	"chat-node/database"
	"chat-node/database/conversations"
	"chat-node/handler"
	"context"
	"fmt"

	"github.com/livekit/protocol/livekit"
)

func join(message handler.Message) {

	if message.ValidateForm("id", "token") {
		handler.ErrorResponse(message, "invalid")
		return
	}

	id := int64(message.Data["id"].(float64)) // ID of the conversation
	token := message.Data["token"].(string)
	claims, valid := calls.GetCallClaims(token)

	if !valid {
		handler.ErrorResponse(message, "invalid")
		return
	}

	// Check if room name is valid
	if !claims.Valid(fmt.Sprintf("c_%d", id)) {
		handler.ErrorResponse(message, "invalid")
		return
	}

	// Check if user is member of the conversation
	if database.DBConn.Where("conversation = ? AND account = ?", id, message.Client.ID).Find(&conversations.Member{}).Error != nil {
		handler.ErrorResponse(message, "invalid")
		return
	}

	// Check if there is already a call (livekit room)
	res, _ := calls.RoomClient.ListRooms(context.Background(), &livekit.ListRoomsRequest{
		Names: []string{claims.CID},
	})

	if len(res.Rooms) > 0 {

		// Connect to call
		tk, err := calls.GetJoinToken(claims.CID, fmt.Sprintf("%d", message.Client.ID))

		if err != nil {
			handler.ErrorResponse(message, "server.error")
			return
		}

		handler.NormalResponse(message, map[string]interface{}{
			"success": true,
			"token":   tk,
		})
	}

	// Create call (livekit room)
	_, err := calls.RoomClient.CreateRoom(context.Background(), &livekit.CreateRoomRequest{
		Name:         claims.CID,
		EmptyTimeout: 5,
	})

	if err != nil {
		handler.ErrorResponse(message, "server.error")
		return
	}

	// Connect to call
	tk, err := calls.GetJoinToken(claims.CID, fmt.Sprintf("%d", message.Client.ID))

	if err != nil {
		handler.ErrorResponse(message, "server.error")
		return
	}

	handler.NormalResponse(message, map[string]interface{}{
		"success": true,
		"token":   tk,
	})
}
