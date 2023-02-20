package ping

import (
	"chat-node/util"

	"github.com/gofiber/fiber/v2"
)

func Pong(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"gateway": util.NODE_ID,
		"app":     "fj.chat-node",
	})
}
