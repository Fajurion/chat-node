package account_routes

import (
	integration "fajurion.com/node-integration"
	"github.com/gofiber/fiber/v2"
)

type StatusSetRequest struct {
	Status      string   `json:"status"` // Encrypted data
	Identifiers []string `json:"identifiers"`
}

func (s *StatusSetRequest) Validate() bool {
	return len(s.Identifiers) < 100 && len(s.Status) > 0
}

// Route: /account/status
func setStatus(c *fiber.Ctx) error {

	// Parse request
	var req StatusSetRequest
	if err := c.BodyParser(&req); err != nil {
		return integration.InvalidRequest(c)
	}

	// Validate request
	if !req.Validate() {
		return integration.InvalidRequest(c)
	}

	return nil
}