package send

import (
	"chat-node/pipe"
	"chat-node/util"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
)

func Socketless(nodeEntity pipe.Node, message pipe.Message) error {

	msg, err := sonic.Marshal(message)
	if err != nil {
		return err
	}

	_, err = util.PostRaw(util.Protocol+nodeEntity.Domain+"/adoption/socketless", fiber.Map{
		"this":    util.NODE_ID,
		"token":   nodeEntity.Token,
		"message": string(msg),
	})

	return err
}
