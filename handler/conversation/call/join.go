package call

import (
	"chat-node/calls"
	"chat-node/database"
	"chat-node/database/conversations"
	"chat-node/handler"
	"chat-node/pipe"
	"chat-node/pipe/send"
	"chat-node/util/requests"
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

	// Check if joiner is the same as the creator
	if claims.Ow == message.Client.ID {
		handler.ErrorResponse(message, "no.join")
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

	// Send to the conversation
	members, nodes, err := requests.LoadConversationDetails(uint(id))
	if err != nil {
		handler.ErrorResponse(message, "server.error")
		return
	}

	// Tell others about it
	send.Pipe(pipe.Message{
		Event: pipe.Event{
			Name: "c_o:l",
			Data: map[string]interface{}{
				"conv": id,
			},
		},
		Channel: pipe.Conversation(members, nodes),
	})

	print("sending to", claims.Ow)

	// Let the owner join
	send.Pipe(pipe.Message{
		Event: pipe.Event{
			Name: "c_s:l",
			Data: map[string]interface{}{
				"conv": id,
			},
		},
		Channel: pipe.BroadcastChannel([]int64{claims.Ow}),
	})

	handler.NormalResponse(message, map[string]interface{}{
		"success": true,
		"token":   tk,
	})
}
