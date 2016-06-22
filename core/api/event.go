package api

import "time"

type Event struct {
	EventType string    `json:"type"`
	Time      time.Time `json:"time"`
}
