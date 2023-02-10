package adoption

import (
	"chat-node/pipe"
	"chat-node/util"

	"github.com/gofiber/fiber/v2"
)

type initializeRequest struct {
	Token string    `json:"token"`
	Node  pipe.Node `json:"node"`
}

func initialize(c *fiber.Ctx) error {

	// Parse request
	var req initializeRequest
	if err := c.BodyParser(&req); err != nil {
		return err
	}

	// Check token
	if req.Token != util.NODE_TOKEN {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	// Add node
	pipe.Nodes[req.Node.ID] = req.Node

	return c.JSON(fiber.Map{
		"success": true,
	})
}
