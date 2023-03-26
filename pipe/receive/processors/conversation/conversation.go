package conversation

import "chat-node/pipe/receive/processors"

func SetupProcessors() {
	processors.Processors["conv_open:l"] = open
}
