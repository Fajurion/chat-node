package routes

import (
	"chat-node/routes/auth"
	"chat-node/routes/ping"

	integration "fajurion.com/node-integration"
	"github.com/Fajurion/pipes"
	"github.com/Fajurion/pipesfiber"
	pipesfutil "github.com/Fajurion/pipesfiber/util"
	"github.com/gofiber/fiber/v2"
)

func Setup(router fiber.Router) {
	router.Route("/auth", auth.Setup)
	router.Post("/ping", ping.Pong)

	pipesfiber.Setup(pipesfutil.Config{
		ExpectedConnections: 10_0_0_0, // 10 thousand, but funny
		NodeDisconnectHandler: func(node pipes.Node) {
			integration.ReportOffline(node)
		},
	})
	router.Route("/", pipesfiber.SetupRoutes)
}
