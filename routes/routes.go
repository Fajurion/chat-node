package routes

import (
	"chat-node/routes/auth"
	conversation_routes "chat-node/routes/conversations"
	mailbox_routes "chat-node/routes/mailbox"
	"chat-node/routes/ping"
	"chat-node/service"
	"chat-node/util"
	"chat-node/util/requests"
	"log"
	"time"

	integration "fajurion.com/node-integration"
	"github.com/Fajurion/pipes"
	"github.com/Fajurion/pipesfiber"
	pipesfroutes "github.com/Fajurion/pipesfiber/routes"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
)

func Setup(router fiber.Router) {

	// Unauthorized routes (for backend/nodes only)
	router.Route("/auth", auth.Setup)
	router.Post("/ping", ping.Pong)

	// Authorized routes (for accounts with remote id only)
	router.Route("/conversations", conversation_routes.SetupRoutes)
	router.Route("/mailbox", mailbox_routes.SetupRoutes)

	// Authorized by using a remote id or normal token
	router.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{
			JWTAlg: jwtware.HS256,
			Key:    []byte(integration.JwtSecret),
		},

		// Checks if the token is expired
		SuccessHandler: func(c *fiber.Ctx) error {

			if util.IsExpired(c) {
				return requests.InvalidRequest(c)
			}

			if !util.IsRemoteId(c) {
				return requests.InvalidRequest(c)
			}

			return c.Next()
		},

		// Error handler
		ErrorHandler: func(c *fiber.Ctx, err error) error {

			log.Println(c.Route().Path, "jwt error:", err.Error())

			// Return error message
			return c.SendStatus(fiber.StatusUnauthorized)
		},
	}))

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

			util.PostRequest("/node/disconnect", map[string]interface{}{
				"node":    util.NODE_ID,
				"token":   util.NODE_TOKEN,
				"session": client.Session,
			})
		},

		// Handle enter network
		ClientConnectHandler: func(client *pipesfiber.Client) bool {
			if integration.Testing {
				log.Println("Client connected:", client.ID)
			}
			disconnect := !service.User(client)
			log.Println("Setup finish")
			if disconnect {
				log.Println("Something went wrong with setup: ", client.ID)
			}
			return disconnect
		},

		// Handle client entering network
		ClientEnterNetworkHandler: func(client *pipesfiber.Client) bool {
			return false
		},
	})
	router.Route("/", pipesfroutes.SetupRoutes)
}
