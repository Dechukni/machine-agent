package op

import "time"

// The base event for all the events published to the client
type Event struct {
	EventType string    `json:"type"`
	Time      time.Time `json:"time"`
}

// The error events, if any error occurs during operation Call processing
type ErrorEvent struct {
	Event
	Message string `json:"message"`
}

// Creates a new error event from the event
func NewErrorEvent(err error) ErrorEvent {
	return ErrorEvent{
		Event: Event{
			EventType: "error",
			Time:      time.Now(),
		},
		Message: err.Error(),
	}
}
