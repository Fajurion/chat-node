package routes

import (
	integration "fajurion.com/node-integration"
	"github.com/Fajurion/pipes"
	"github.com/Fajurion/pipes/receive"
	"github.com/Fajurion/pipes/send"
	"github.com/gofiber/fiber/v2"
)

type socketlessEvent struct {
	Token   string        `json:"token"`
	Message pipes.Message `json:"message"`
}

func socketless(c *fiber.Ctx) error {

	// Parse request
	var event socketlessEvent
	if err := integration.BodyParser(c, &event); err != nil {
		return err
	}

	// Check token
	if event.Token != pipes.CurrentNode.Token {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	receive.HandleMessage(send.ProtocolWS, event.Message)

	return integration.SuccessfulRequest(c)
}
