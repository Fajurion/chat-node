package routes

import (
	"chat-node/database/fetching"
	"chat-node/handler/account"
	"chat-node/routes/auth"
	"chat-node/routes/ping"
	"chat-node/service"
	"log"
	"time"

	integration "fajurion.com/node-integration"
	"github.com/Fajurion/pipes"
	"github.com/Fajurion/pipesfiber"
	pipesfroutes "github.com/Fajurion/pipesfiber/routes"
	"github.com/gofiber/fiber/v2"
)

func Setup(router fiber.Router) {
	router.Route("/auth", auth.Setup)
	router.Post("/ping", ping.Pong)

	pipesfiber.Setup(pipesfiber.Config{
		ExpectedConnections: 10_0_0_0,       // 10 thousand, but funny
		SessionDuration:     time.Hour * 24, // This is kinda important

		// Report nodes as offline
		NodeDisconnectHandler: func(node pipes.Node) {
			integration.ReportOffline(node)
		},

		// Handle client disconnect
		ClientDisconnectHandler: func(client *pipesfiber.Client) {
			if integration.Testing {
				log.Println("Client disconnected:", client.ID)
			}
			account.UpdateStatus(client, fetching.StatusOffline, "", false)
		},

		// Handle enter network
		ClientConnectHandler: func(client *pipesfiber.Client) bool {
			if integration.Testing {
				log.Println("Client connected:", client.ID)
			}
			disconnect := !service.User(client)
			if disconnect {
				log.Println("Something went wrong with setup: ", client.ID)
			}
			return disconnect
		},
	})
	router.Route("/", pipesfroutes.SetupRoutes)
}
