package routes

import (
	account_routes "chat-node/routes/account"
	"chat-node/routes/auth"
	conversation_routes "chat-node/routes/conversations"
	"chat-node/routes/ping"
	"chat-node/service"
	"chat-node/util"
	"encoding/base64"
	"log"
	"time"

	integration "fajurion.com/node-integration"
	"github.com/Fajurion/pipes"
	"github.com/Fajurion/pipes/adapter"
	"github.com/Fajurion/pipesfiber"
	pipesfroutes "github.com/Fajurion/pipesfiber/routes"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
)

func Setup(router fiber.Router) {

	// Return the public key for TC Protection
	router.Post("/pub", func(c *fiber.Ctx) error {

		// Return the public key in a packaged form (string)
		return c.JSON(fiber.Map{
			"pub": integration.PackageRSAPublicKey(integration.NodePublicKey),
		})
	})

	router.Post("/ping", ping.Pong)

	router.Route("/", encryptedRoutes)
}

func encryptedRoutes(router fiber.Router) {

	// Add Through Cloudflare Protection middleware
	router.Use(func(c *fiber.Ctx) error {

		// Get the AES encryption key from the Auth-Tag header
		aesKeyEncoded, valid := c.GetReqHeaders()["Auth-Tag"]
		if !valid {
			log.Println("no header")
			return c.SendStatus(fiber.StatusPreconditionFailed)
		}
		aesKeyEncrypted, err := base64.StdEncoding.DecodeString(aesKeyEncoded)
		if err != nil {
			log.Println("no decoding")
			return c.SendStatus(fiber.StatusPreconditionFailed)
		}

		// Decrypt the AES key using the private key of this node
		aesKey, err := integration.DecryptRSA(integration.NodePrivateKey, aesKeyEncrypted)
		if err != nil {
			return c.SendStatus(fiber.StatusPreconditionRequired)
		}

		// Decrypt the request body using the key attached to the Auth-Tag header
		decrypted, err := integration.DecryptAES(aesKey, c.Body())
		if err != nil {
			return c.SendStatus(fiber.StatusNetworkAuthenticationRequired)
		}

		// Set some variables for use when sending back the response
		c.Locals("body", decrypted)
		c.Locals("key", aesKey)
		c.Locals("srv_pub", integration.NodePrivateKey)

		// Go to the next middleware/handler
		return c.Next()
	})

	// Unauthorized routes (for backend/nodes only)
	router.Route("/auth", auth.Setup)

	setupPipesFiber(router)

	// Authorized by using a remote id or normal token
	router.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{
			JWTAlg: jwtware.HS256,
			Key:    []byte(integration.JwtSecret),
		},

		// Checks if the token is expired
		SuccessHandler: func(c *fiber.Ctx) error {

			// Check if the JWT is expired
			if util.IsExpired(c) {
				return integration.InvalidRequest(c, "expired remote id")
			}

			// Check if the JWT is a remote id
			if !util.IsRemoteId(c) {
				return integration.InvalidRequest(c, "jwt isn't a remote id")
			}

			// Go to the next middleware/handler
			return c.Next()
		},

		// Error handler
		ErrorHandler: func(c *fiber.Ctx, err error) error {

			log.Println(c.Route().Path, "jwt error:", err.Error())

			// Return error message
			return c.SendStatus(fiber.StatusUnauthorized)
		},
	}))

	// Authorized routes (for accounts with remote id only)
	router.Route("/conversations", conversation_routes.SetupRoutes)
	router.Route("/account", account_routes.SetupRoutes)
}

func setupPipesFiber(router fiber.Router) {
	adapter.SetupCaching()
	pipesfiber.Setup(pipesfiber.Config{
		ExpectedConnections: 10_0_0_0,       // 10 thousand, but funny
		SessionDuration:     time.Hour * 24, // This is kinda important

		// Report nodes as offline
		NodeDisconnectHandler: func(node pipes.Node) {

			// Report that a node is offline to the backend
			integration.ReportOffline(node)
		},

		// Handle client disconnect
		ClientDisconnectHandler: func(client *pipesfiber.Client) {

			// Print debug stuff if in debug mode
			if integration.Testing {
				log.Println("Client disconnected:", client.ID)
			}

			// Tell the backend that someone disconnected
			integration.PostRequest("/node/disconnect", map[string]interface{}{
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

			// Initialize the user and check if he needs to be disconnected
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
