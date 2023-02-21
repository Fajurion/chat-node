package send

import (
	"chat-node/bridge"
	"chat-node/pipe"

	"github.com/bytedance/sonic"
)

func Client(id int64, event pipe.Event) {

	msg, err := sonic.Marshal(event)
	if err != nil {
		return
	}

	bridge.Send(id, msg)
}
