package conversation

import (
	"github.com/Fajurion/pipes/receive/processors"
)

func SetupProcessors() {
	processors.Processors["conv_open:l"] = open
}
