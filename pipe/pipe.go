package pipe

type Event struct {
	Name string                 `json:"name"`
	Data map[string]interface{} `json:"data"`
}

type Channel struct {
	Sender  int64   `json:"sender"`  // User ID (0 for system)
	Channel string  `json:"channel"` // "project", "event"
	Target  []int64 `json:"target"`  // Project ID or User ID
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

func (c Channel) IsSocketless() bool {
	return c.Channel == "socketless"
}

func P2PChannel(sender int64, receiver int64, receiverNode int64) Channel {
	return Channel{
		Channel: "p2p",
		Sender:  sender,
		Target:  []int64{receiver, receiverNode},
	}
}

func ProjectChannel(sender int64, project int64) Channel {
	return Channel{
		Channel: "project",
		Sender:  sender,
		Target:  []int64{project},
	}
}

func BroadcastChannel(sender int64, receivers []int64) Channel {
	return Channel{
		Channel: "broadcast",
		Sender:  sender,
		Target:  receivers,
	}
}
