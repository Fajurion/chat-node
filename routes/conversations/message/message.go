package message_routes

import (
	"chat-node/caching"
	"chat-node/database"
	"chat-node/database/conversations"
	"chat-node/util"
	"time"

	"github.com/Fajurion/pipes"
	"github.com/Fajurion/pipes/send"
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
)

const systemSender = "6969"

func SetupRoutes(router fiber.Router) {
	router.Post("/send", sendMessage)
	router.Post("/delete", deleteMessage)
}

// System messages
const DeletedConversation = "conv.deleted"
const DeletedMessage = "msg.deleted"
const GroupNewAdmin = "group.new_admin"
const GroupRankChange = "group.rank_change"
const GroupMemberJoin = "group.member_join"
const GroupMemberLeave = "group.member_leave"

func SendSystemMessage(conversation string, content string, attachments []string) error {

	contentJson, err := sonic.MarshalString(map[string]interface{}{
		"c": content,
		"a": attachments,
	})
	if err != nil {
		return err
	}

	messageId := util.GenerateToken(32)
	message := conversations.Message{
		ID:           messageId,
		Conversation: conversation,
		Certificate:  "",
		Data:         contentJson,
		Sender:       systemSender,
		Creation:     time.Now().UnixMilli(),
		Edited:       false,
	}

	if err := database.DBConn.Create(&message).Error; err != nil {
		return err
	}

	// Load members
	members, err := caching.LoadMembers(conversation)
	if err != nil {
		return err
	}
	adapters, nodes := caching.MembersToPipes(members)

	event := MessageEvent(message)
	err = send.Pipe(send.ProtocolWS, pipes.Message{
		Channel: pipes.Conversation(adapters, nodes),
		Event:   event,
	})
	if err != nil {
		return err
	}

	return nil
}

func AttachAccount(encrypted string) string {
	return "a:" + encrypted
}
