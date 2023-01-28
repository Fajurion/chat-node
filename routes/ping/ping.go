package ping

import (
	"chat-node/util"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func Pong(c *fiber.Ctx) error {
	return c.SendString(fmt.Sprintf("Node %d | Fajurion Network", util.NODE_ID))
}
