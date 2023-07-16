package call

import (
	"chat-node/calls"
	"chat-node/database"
	"chat-node/database/conversations"
	"context"
	"fmt"

	"github.com/Fajurion/pipesfiber/wshandler"
	"github.com/livekit/protocol/livekit"
)

func status(message wshandler.Message) {

	if message.ValidateForm("id") {
		wshandler.ErrorResponse(message, "invalid")
		return
	}

	// Check if user is member
	if database.DBConn.Model(&conversations.Member{}).Where("conversation = ?", message.Data["id"]).Error != nil {
		wshandler.ErrorResponse(message, "invalid")
		return
	}

	// Check for call
	members, err := calls.RoomClient.ListParticipants(context.Background(), &livekit.ListParticipantsRequest{
		Room: fmt.Sprintf("c_%d", int64(message.Data["id"].(float64))),
	})

	if err != nil {
		wshandler.NormalResponse(message, map[string]interface{}{
			"success": true,
			"call":    false,
		})
		return
	}

	var participants []string
	for _, participant := range members.Participants {
		participants = append(participants, participant.Identity)
	}

	wshandler.NormalResponse(message, map[string]interface{}{
		"success": true,
		"call":    true,
		"members": participants,
	})
}
