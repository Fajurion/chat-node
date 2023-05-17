package adoption

import (
	"chat-node/util"

	"github.com/Fajurion/pipes"
	"github.com/Fajurion/pipes/receive"
	"github.com/Fajurion/pipes/send"
	"github.com/gofiber/fiber/v2"
)

type socketLessEvent struct {
	Token   string        `json:"token"`
	This    uint          `json:"this"`
	Message pipes.Message `json:"message"`
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

	receive.HandleMessage(send.ProtocolWS, event.Message)

	return c.JSON(fiber.Map{
		"success": true,
	})
}
