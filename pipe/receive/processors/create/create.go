package create

import (
	"chat-node/pipe/receive/processors/conversation"
)

func SetupProcessors() {

	// Call the setup functions
	conversation.SetupProcessors()
}
