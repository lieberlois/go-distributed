package coordinator

import (
	"time"
)

type EventAggregator struct {
	listeners map[string][]func(EventData)
}

type EventData struct {
	Name      string
	Value     float64
	Timestamp time.Time
}

func NewEventAggregator() *EventAggregator {
	return &EventAggregator{listeners: make(map[string][]func(EventData))}
}

func (ea *EventAggregator) AddListener(name string, f func(EventData)) {
	ea.listeners[name] = append(ea.listeners[name], f)
}

func (ea *EventAggregator) PublishEvent(name string, eventData EventData) {
	if ea.listeners[name] != nil {
		for _, cb := range ea.listeners[name] {
			cb(eventData)
		}
	}
}