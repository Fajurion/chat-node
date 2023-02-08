package pipe

import (
	"context"
	"fmt"
	"log"
	"time"

	"nhooyr.io/websocket"
)

var NodeConnections map[int64]*websocket.Conn

func ConnectToNode(node Node) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	// Connect to node
	c, _, err := websocket.Dial(ctx, node.GetWebSocket(), &websocket.DialOptions{
		Subprotocols: []string{fmt.Sprintf("%s_%d_%s", node.Token, node.App, node.Domain)},
	})

	if err != nil {
		return
	}

	// Add connection to map
	NodeConnections[node.ID] = c

	log.Printf("Outgoing event stream to node %d connected.", node.ID)
}
