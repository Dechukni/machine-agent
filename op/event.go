package op

import "time"

const (
	ErrorEventType = "error"
)

// The base event for all the events published to the client
type Event struct {
	EventType string      `json:"type"`
	Time      time.Time   `json:"time"`
	Body      interface{} `json:"body"`
}

// Creates a new error event from the Error
func NewErrorEvent(err Error) *Event {
	return NewEventNow(ErrorEventType, err)
}

func NewEventNow(eType string, Body interface{}) *Event {
	return &Event{
		EventType: eType,
		Time:      time.Now(),
		Body:      Body,
	}
}

func NewEvent(eType string, Body interface{}, time time.Time) *Event {
	return &Event{
		EventType: eType,
		Time:      time,
		Body:      Body,
	}
}
