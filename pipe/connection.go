package pipe

import (
	"context"
	"log"
	"time"

	"nhooyr.io/websocket"
)

var NodeConnections map[int64]*websocket.Conn

func ConnectToNode(node Node) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	// Connect to node
	c, _, err := websocket.Dial(ctx, node.GetWebSocket(), nil)

	if err != nil {
		return
	}

	// Add connection to map
	NodeConnections[node.ID] = c

	log.Println("Connected to node ", node.Domain, ".")
}
