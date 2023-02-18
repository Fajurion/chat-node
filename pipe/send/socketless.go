package send

import (
	"chat-node/pipe"
	"chat-node/util"

	"github.com/gofiber/fiber/v2"
)

func sendSocketless(message pipe.Message, msg []byte) error {

	_, err := util.PostRaw(util.Protocol+message.Channel.Node.Domain+"/adoption/socketless", fiber.Map{
		"this":    util.NODE_ID,
		"token":   message.Channel.Node.Token,
		"message": string(msg),
	})

	return err
}
