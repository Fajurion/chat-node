package message_routes

import (
	"fmt"

	integration "fajurion.com/node-integration"
	"github.com/gofiber/fiber/v2"
)

type sendSystemRequest struct {
	Conversation string
}

func sendSystem(c *fiber.Ctx) error {

	var req sendSystemRequest
	if integration.BodyParser(c, &req) != nil {
		return integration.InvalidRequest(c, "invalid request")
	}

	err := SendSystemMessage(req.Conversation, "group.rank_change", []string{"1", "2", "DtLmwVF35oiE", "NZJNP232RS5g"})
	if err != nil {
		return integration.InvalidRequest(c, fmt.Sprintf("couldn't send system message: %s", err.Error()))
	}

	return integration.SuccessfulRequest(c)
}
