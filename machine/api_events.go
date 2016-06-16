// TODO add subscribed event types
package machine

import "time"

const (
	CONNECTED = "connected"

	PROCESS_STARTED = "process_started"
	PROCESS_DIED    = "process_died"
	STDOUT          = "stdout"
	STDERR          = "stderr"
)

type Event struct {
	EventType string    `json:"type"`
	Time      time.Time `json:"time"`
}

type ProcessEvent struct {
	Event
	Pid uint64 `json:"pid"`
}

type ProcessStatusEvent struct {
	ProcessEvent
	NativePid   uint64 `json:"nativePid"`
	Name        string `json:"name"`
	CommandLine string `json:"commandLine"`
}

type ProcessOutputEvent struct {
	ProcessEvent
	Text string `json:"text"`
}

type ChannelEvent struct {
	Event
	ChannelId string `json:"channelId"`
	Text      string `json:"text"`
}
