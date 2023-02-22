package send

import (
	"chat-node/pipe"
	"chat-node/util"

	"github.com/gofiber/fiber/v2"
)

func Socketless(nodeEntity pipe.Node, message pipe.Message) error {

	_, err := util.PostRaw(util.Protocol+nodeEntity.Domain+"/adoption/socketless", fiber.Map{
		"this":    util.NODE_ID,
		"token":   nodeEntity.Token,
		"message": message,
	})

	return err
}
