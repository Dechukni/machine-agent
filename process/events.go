// TODO add subscribed event types
package process

import "github.com/evoevodin/machine-agent/op"

const (
	PROCESS_STARTED = "process_started"
	PROCESS_DIED    = "process_died"
	STDOUT          = "stdout"
	STDERR          = "stderr"
)

type ProcessEvent struct {
	op.Event
	Pid uint64 `json:"pid"`
}

type ProcessStatusEvent struct {
	ProcessEvent
	NativePid   int    `json:"nativePid"`
	Name        string `json:"name"`
	CommandLine string `json:"commandLine"`
}

type ProcessOutputEvent struct {
	ProcessEvent
	Text string `json:"text"`
}
