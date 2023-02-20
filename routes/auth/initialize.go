package auth

import (
	"chat-node/bridge"
	"chat-node/util"
	"chat-node/util/requests"

	"github.com/gofiber/fiber/v2"
)

type intializeRequest struct {
	NodeToken string `json:"node_token"`
	Session   uint64 `json:"session"`
	UserID    int64  `json:"user_id"`
	Username  string `json:"username"`
	Tag       string `json:"tag"`
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

	// Check if there are too many users
	if bridge.GetConnections(req.UserID) >= 3 {
		return requests.FailedRequest(c, "too.many.connections", nil)
	}

	bridge.AddToken(tk, req.UserID, req.Session, req.Username, req.Tag)

	return c.JSON(fiber.Map{
		"success": true,
		"load":    0,
		"token":   tk,
	})
}
