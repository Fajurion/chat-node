package auth

import "github.com/gofiber/fiber/v2"

func Setup(router fiber.Router) {
	router.Post("/initialize", initializeConnection)
}
