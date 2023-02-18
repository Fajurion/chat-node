package adoption

import (
	"chat-node/pipe"
	"chat-node/pipe/receive"
	"chat-node/util"

	"github.com/gofiber/fiber/v2"
)

type socketLessEvent struct {
	Token   string       `json:"token"`
	This    uint         `json:"this"`
	Message pipe.Message `json:"message"`
}

func socketless(c *fiber.Ctx) error {

	// Parse request
	var event socketLessEvent
	if err := c.BodyParser(&event); err != nil {
		return err
	}

	// Check token
	if event.Token != util.NODE_TOKEN {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	receive.Handle(event.Message)

	return c.JSON(fiber.Map{
		"success": true,
	})
}
