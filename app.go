package main

import (
	"chat-node/caching"
	"chat-node/calls"
	"chat-node/database"
	"chat-node/handler"
	"chat-node/processors"
	"chat-node/routes"
	"chat-node/util"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	integration "fajurion.com/node-integration"
	"github.com/Fajurion/pipes"
	"github.com/Fajurion/pipes/connection"
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

var APP_ID uint = 0

func main() {

	fmt.Println("IF YOU ARE ON LINUX, MAKE SURE TO RUN THIS PROGRAM WITH RIGHT PERMISSIONS TO NODE_ENV")
	log.SetOutput(os.Stdout)

	// Setting up the node
	if !integration.Setup() {
		return
	}

	// Connect to the database
	database.Connect()
	caching.SetupCaches()

	// Create fiber app
	app := fiber.New(fiber.Config{
		JSONEncoder: sonic.Marshal,
		JSONDecoder: sonic.Unmarshal,
	})

	pipes.SetupCurrent(fmt.Sprintf("%d", integration.NODE_ID), integration.NODE_TOKEN)

	// Query current node
	_, tk, currentApp, domain := integration.GetCurrent()
	util.NODE_TOKEN = tk
	APP_ID = currentApp

	// Report online status
	res := integration.SetOnline()
	parseNodes(res)

	pipes.SetupSocketless(domain + "/adoption/socketless")

	app.Use(logger.New())
	app.Route("/", routes.Setup)

	// Connect to livekit
	calls.Connect()

	// Create handlers
	handler.Create()

	// Initialize processors
	processors.SetupProcessors()

	// Check if test mode or production
	args := strings.Split(domain, ":")
	var port int
	var err error
	if os.Getenv("OVERWRITE_PORT") != "" {
		port, err = strconv.Atoi(os.Getenv("OVERWRITE_PORT"))
	} else {
		port, err = strconv.Atoi(args[1])
	}
	if err != nil {
		log.Println("Error: Couldn't parse port of current node")
		return
	}

	protocol := os.Getenv("WEBSOCKET_PROTOCOL")
	if protocol == "" {
		protocol = "wss://"
	}
	pipes.SetupWS(protocol + domain + "/connect")

	// Connect to other nodes
	pipes.IterateNodes(func(_ string, node pipes.Node) bool {

		log.Println("Connecting to node " + node.WS)

		if err := connection.ConnectWS(node); err != nil {
			log.Println(err.Error())
		}
		return true
	})

	pipes.DebugLogs = true // TODO: Replace in production
	if integration.Testing {

		// Start on localhost
		app.Listen(fmt.Sprintf("localhost:%d", port))
	} else {

		// Start on all interfaces
		app.Listen(fmt.Sprintf("0.0.0.0:%d", port))
	}
}

// Shared function between all nodes
func parseNodes(res map[string]interface{}) bool {

	if res["nodes"] == nil {
		return true
	}

	nodeList := res["nodes"].([]interface{})

	for _, node := range nodeList {
		n := node.(map[string]interface{})

		// Extract port and domain
		args := strings.Split(n["domain"].(string), ":")
		domain := args[0]
		port, err := strconv.Atoi(args[1])
		if err != nil {
			log.Println("Error: Couldn't parse port of node " + n["id"].(string))
			return true
		}

		// Add node to pipes
		pipes.AddNode(pipes.Node{
			ID:    fmt.Sprintf("%d", int64(n["id"].(float64))),
			Token: n["token"].(string),
			WS:    "ws://" + fmt.Sprintf("%s:%d", domain, port) + "/adoption",
		})
	}

	return false
}
