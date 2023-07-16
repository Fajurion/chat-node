package call

import (
	"chat-node/calls"
	"chat-node/database"
	"chat-node/database/conversations"
	"chat-node/util/requests"
	"context"

	"github.com/Fajurion/pipes"
	"github.com/Fajurion/pipes/send"
	"github.com/Fajurion/pipesfiber/wshandler"
	"github.com/livekit/protocol/livekit"
)

func join(message wshandler.Message) {

	if message.ValidateForm("id", "token") {
		wshandler.ErrorResponse(message, "invalid")
		return
	}

	id := message.Data["id"].(string) // ID of the conversation
	token := message.Data["token"].(string)
	claims, valid := calls.GetCallClaims(token)

	if !valid {
		wshandler.ErrorResponse(message, "invalid")
		return
	}

	// Check if joiner is the same as the creator
	if claims.Ow == message.Client.ID {
		wshandler.ErrorResponse(message, "no.join")
		return
	}

	// Check if room name is valid
	if !claims.Valid(id) {
		wshandler.ErrorResponse(message, "invalid")
		return
	}

	// Check if user is member of the conversation
	if database.DBConn.Where("conversation = ? AND account = ?", id, message.Client.ID).Find(&conversations.Member{}).Error != nil {
		wshandler.ErrorResponse(message, "invalid")
		return
	}

	// Check if there is already a call (livekit room)
	res, _ := calls.RoomClient.ListRooms(context.Background(), &livekit.ListRoomsRequest{
		Names: []string{claims.CID},
	})

	if len(res.Rooms) > 0 {

		// Connect to call
		tk, err := calls.GetJoinToken(claims.CID, message.Client.ID)

		if err != nil {
			wshandler.ErrorResponse(message, "server.error")
			return
		}

		wshandler.NormalResponse(message, map[string]interface{}{
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
		wshandler.ErrorResponse(message, "server.error")
		return
	}

	// Connect to call
	tk, err := calls.GetJoinToken(claims.CID, message.Client.ID)

	if err != nil {
		wshandler.ErrorResponse(message, "server.error")
		return
	}

	// Send to the conversation
	members, nodes, err := requests.LoadConversationDetails(id)
	if err != nil {
		wshandler.ErrorResponse(message, "server.error")
		return
	}

	// Tell others about it
	send.Pipe(send.ProtocolWS, pipes.Message{
		Event: pipes.Event{
			Name: "c_o:l",
			Data: map[string]interface{}{
				"conv": id,
			},
		},
		Channel: pipes.Conversation(members, nodes),
	})

	print("sending to", claims.Ow)

	// Let the owner join
	send.Pipe(send.ProtocolWS, pipes.Message{
		Event: pipes.Event{
			Name: "c_s:l",
			Data: map[string]interface{}{
				"conv": id,
			},
		},
		Channel: pipes.BroadcastChannel([]string{claims.Ow}),
	})

	wshandler.NormalResponse(message, map[string]interface{}{
		"success": true,
		"token":   tk,
	})
}
