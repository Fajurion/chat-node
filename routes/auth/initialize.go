package auth

import (
	"chat-node/util"

	integration "fajurion.com/node-integration"
	"github.com/Fajurion/pipesfiber"
	"github.com/gofiber/fiber/v2"
)

const SenderUser = 0
const SenderNode = 1

type intializeRequest struct {
	Sender    uint   `json:"sender"`
	Account   string `json:"account"`
	Session   string `json:"session"`
	NodeToken string `json:"node_token"`
	End       int64  `json:"end"`
}

func initializeConnection(c *fiber.Ctx) error {

	// Parse the request
	var req intializeRequest
	if err := integration.BodyParser(c, &req); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if req.Sender == SenderNode {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if util.NODE_TOKEN != req.NodeToken {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	tk := util.GenerateToken(200)

	// Check if there are too many users
	if pipesfiber.GetConnections(req.Account) >= 3 {
		return integration.FailedRequest(c, "too.many.connections", nil)
	}

	pipesfiber.AddToken(tk, pipesfiber.ConnectionToken{
		UserID:  req.Account,
		Session: req.Session,
		Data:    nil,
	})

	return integration.ReturnJSON(c, fiber.Map{
		"success": true,
		"load":    0,
		"token":   tk,
	})
}
