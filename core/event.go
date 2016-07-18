package core

import "time"

type Event struct {
	EventType string    `json:"type"`
	Time      time.Time `json:"time"`
}

type ErrorEvent struct {
	Event
	Message string
}

func NewErrorEvent(err error) ErrorEvent {
	return ErrorEvent{
		Event{
			"error",
			time.Now(),
		},
		err.Error(),
	}
}
