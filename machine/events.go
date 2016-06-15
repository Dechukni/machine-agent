package machine

import "time"

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

// 2016-06-15T20:29:44.437650129+03:00
