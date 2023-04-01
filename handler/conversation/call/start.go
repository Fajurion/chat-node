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

func start(message handler.Message) {

	if message.ValidateForm("id") {
		handler.ErrorResponse(message, "invalid")
		return
	}

	id := int64(message.Data["id"].(float64)) // ID of the conversation
	roomName := fmt.Sprintf("c_%d", id)       // Livekit room name

	// Check if user is member of the conversation
	if database.DBConn.Where("conversation = ? AND account = ?", id, message.Client.ID).Find(&conversations.Member{}).Error != nil {
		handler.ErrorResponse(message, "invalid")
		return
	}

	// Check if there is already a call (livekit room)
	res, _ := calls.RoomClient.ListRooms(context.Background(), &livekit.ListRoomsRequest{
		Names: []string{roomName},
	})

	if len(res.Rooms) > 0 {

		// Connect to call
		tk, err := calls.GetJoinToken(roomName, fmt.Sprintf("%d", message.Client.ID))

		if err != nil {
			handler.ErrorResponse(message, "server.error")
			return
		}

		handler.NormalResponse(message, map[string]interface{}{
			"success": true,
			"token":   tk,
		})

		return
	}

	// Create call invite
	token, err := calls.GenerateCallToken(roomName, message.Client.ID)

	if err != nil {
		handler.ErrorResponse(message, "server.error")
		return
	}

	handler.NormalResponse(message, map[string]interface{}{
		"success": true,
		"id":      id,
		"token":   token,
	})
}
