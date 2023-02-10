package pipe

import (
	"chat-node/util"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"nhooyr.io/websocket"
)

var NodeConnections map[int64]*websocket.Conn = make(map[int64]*websocket.Conn)

func ConnectToNode(node Node) {

	// Connect to node
	c, _, err := websocket.Dial(context.Background(), node.GetWebSocket(), &websocket.DialOptions{
		Subprotocols: []string{fmt.Sprintf("%s_%d_%s", node.Token, util.NODE_ID, util.NODE_TOKEN)},
	})

	if err != nil {
		return
	}

	// Add connection to map
	NodeConnections[node.ID] = c

	log.Printf("Outgoing event stream to node %d connected.", node.ID)

	go func() {
		for {
			time.Sleep(time.Second * 5)

			// Send ping
			c.Write(context.Background(), websocket.MessageText, []byte("ping"))
		}
	}()
}

func ReportOffline(node Node) {

	// Check if connection exists
	if NodeConnections[node.ID] == nil {
		return
	}

	res, err := util.PostRequest("/node/status/offline", fiber.Map{
		"token": node.Token,
	})

	if err != nil {
		log.Println("Failed to report offline status. Is the backend online?")
	}

	log.Println(res)
}
