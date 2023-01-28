package auth

import (
	"chat-node/bridge"
	"chat-node/util"

	"github.com/gofiber/fiber/v2"
)

type intializeRequest struct {
	NodeToken string `json:"node_token"`
	Session   string `json:"session"`
	UserID    int64  `json:"user_id"`
}

func initializeConnection(c *fiber.Ctx) error {

	// Parse the request
	var req intializeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if util.NODE_TOKEN != req.NodeToken {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	tk := util.GenerateToken(200)
	bridge.AddToken(tk, req.UserID, req.Session)

	return c.JSON(fiber.Map{
		"token": tk,
	})
}
