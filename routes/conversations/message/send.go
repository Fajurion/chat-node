package message_routes

import (
	"chat-node/caching"
	"chat-node/database/conversations"
	"chat-node/util"
	"log"
	"time"

	integration "fajurion.com/node-integration"
	"github.com/Fajurion/pipes"
	"github.com/Fajurion/pipes/send"
	"github.com/gofiber/fiber/v2"
)

type MessageSendRequest struct {
	Conversation string `json:"conversation"`
	TokenID      string `json:"token_id"`
	Token        string `json:"token"`
	Data         string `json:"data"`
}

func (r *MessageSendRequest) Validate() bool {
	return len(r.Conversation) > 0 && len(r.Data) > 0 && len(r.Token) == util.ConversationTokenLength
}

// Route: /conversations/message/send
func sendMessage(c *fiber.Ctx) error {

	var req MessageSendRequest
	if err := c.BodyParser(&req); err != nil {
		return integration.InvalidRequest(c)
	}

	// Validate request
	if !req.Validate() {
		return integration.InvalidRequest(c)
	}

	if conversations.CheckSize(req.Data) {
		return integration.FailedRequest(c, "too.big", nil)
	}

	// Validate conversation token
	_, err := caching.ValidateToken(req.TokenID, req.Token)
	if err != nil {
		return integration.InvalidRequest(c)
	}

	// Load members
	members, err := caching.LoadMembers(req.Conversation)
	if err != nil {
		return integration.FailedRequest(c, "server.error", nil)
	}

	found := false
	for _, member := range members {
		if member.Token == req.Token {
			found = true
		}
	}

	if !found {
		return integration.InvalidRequest(c)
	}

	messageId := util.GenerateToken(32)
	certificate, err := conversations.GenerateCertificate(messageId, req.TokenID)
	if err != nil {
		return integration.FailedRequest(c, "server.error", nil)
	}

	message := conversations.Message{
		ID:           messageId,
		Conversation: req.Conversation,
		Certificate:  certificate,
		Data:         req.Data,
		Sender:       req.TokenID,
		Creation:     time.Now().UnixMilli(),
		Edited:       false,
	}

	// Messages aren't stored for now
	/*
		if err := database.DBConn.Create(&message).Error; err != nil {
			return integration.FailedRequest(c, "server.error", err)
		}
	*/

	adapters, nodes := caching.MembersToPipes(members)

	log.Println(adapters)
	log.Println(nodes)

	event := pipes.Event{
		Sender: send.SenderSystem,
		Name:   "conv_msg",
		Data: map[string]interface{}{
			"conv": req.Conversation,
			"msg":  message,
		},
	}

	log.Println("sending message..")

	send.Pipe(send.ProtocolWS, pipes.Message{
		Channel: pipes.Conversation(adapters, nodes),
		Event:   event,
	})

	return c.JSON(fiber.Map{
		"success": true,
		"message": message,
	})
}
