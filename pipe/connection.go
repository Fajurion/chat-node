package pipe

import (
	"chat-node/util"
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
		Subprotocols: []string{fmt.Sprintf("%s_%d_%s", node.Token, util.NODE_ID, util.NODE_TOKEN)},
	})

	if err != nil {
		return
	}

	// Add connection to map
	NodeConnections[node.ID] = c

	log.Printf("Outgoing event stream to node %d connected.", node.ID)
}
