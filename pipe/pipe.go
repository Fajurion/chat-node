package pipe

type Event struct {
	Sender int64                  `json:"sender"` // User ID (0 for system)
	Name   string                 `json:"name"`
	Data   map[string]interface{} `json:"data"`
}

type Channel struct {
	Channel string  `json:"channel"` // "project", "event"
	Target  []int64 `json:"target"`  // Project ID or User ID
	Nodes   []int64 `json:"-"`       // Nodes to send to (only for conversation channel)
}

type Message struct {
	Channel Channel `json:"channel"`
	Event   Event   `json:"event"`
}

func (c Channel) IsP2P() bool {
	return c.Channel == "p2p"
}

func (c Channel) IsProject() bool {
	return c.Channel == "project"
}

func (c Channel) IsBroadcast() bool {
	return c.Channel == "broadcast"
}

func P2PChannel(receiver int64, receiverNode int64) Channel {
	return Channel{
		Channel: "p2p",
		Target:  []int64{receiver, receiverNode},
	}
}

func Conversation(receivers []int64, nodes []int64) Channel {
	return Channel{
		Channel: "conversation",
		Target:  receivers,
		Nodes:   nodes,
	}
}

func BroadcastChannel(receivers []int64) Channel {
	return Channel{
		Channel: "broadcast",
		Target:  receivers,
	}
}
