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

func (c Channel) IsValid(event Event) bool {
	for _, channel := range AvailableChannels {
		if c.Channel == channel {
			return true
		}
	}

	for _, event := range AvailableEvents[c.Channel] {
		if c.Channel == event {
			return true
		}
	}

	if len(c.Target) == 0 {
		return false
	}

	switch c.Channel {
	case "p2p":
	case "project":
		return len(c.Target) > 1

	case "event":
		return true
	}

	return false
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

var AvailableChannels = []string{"project", "broadcast", "p2p"}

var AvailableEvents = map[string][]string{
	"project":   {"msg", "msg_edit", "msg_delete", "conv_join", "conv_leave"},
	"broadcast": {"status_change"},
	"p2p":       {"key_exc", "friend_rq", "friend_rq_accept", "friend_rq_reject"},
}

func P2PChannel(sender int64, receiver int64) Channel {
	return Channel{
		Channel: "p2p",
		Sender:  sender,
		Target:  []int64{receiver},
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
