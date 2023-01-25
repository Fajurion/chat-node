package ping

import "github.com/gofiber/fiber/v2"

func Pong(c *fiber.Ctx) error {
	return c.SendString("Pong")
}
