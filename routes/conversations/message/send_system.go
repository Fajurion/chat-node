package message_routes

import (
	"chat-node/util/requests"

	"github.com/gofiber/fiber/v2"
)

type sendSystemRequest struct {
	Conversation string
}

func sendSystem(c *fiber.Ctx) error {

	var req sendSystemRequest
	if c.BodyParser(&req) != nil {
		return requests.InvalidRequest(c)
	}

	err := SendSystemMessage(req.Conversation, "group.rank_change", []string{"1", "2", "DtLmwVF35oiE", "NZJNP232RS5g"})
	if err != nil {
		return requests.InvalidRequest(c)
	}

	return requests.SuccessfulRequest(c)
}
