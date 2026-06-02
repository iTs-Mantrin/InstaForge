package events

import (
	"fmt"
	"time"
)

// EventType categorizes domain events.
type EventType string

const (
	EventDownloadRequested EventType = "download.requested"
	EventDownloadStarted   EventType = "download.started"
	EventDownloadCompleted EventType = "download.completed"
	EventDownloadFailed    EventType = "download.failed"

	EventUserSearched  EventType = "user.searched"
	EventURLProcessed  EventType = "url.processed"

	EventRateLimitHit    EventType = "ratelimit.hit"
	EventAbuseDetected   EventType = "abuse.detected"
	EventAPIKeyUsed      EventType = "apikey.used"

	EventCacheHit  EventType = "cache.hit"
	EventCacheMiss EventType = "cache.miss"

	EventWorkerStarted   EventType = "worker.started"
	EventWorkerStopped   EventType = "worker.stopped"
	EventWorkerPanic     EventType = "worker.panic"
)

// Event represents a domain event for observability and event-driven processing.
type Event struct {
	ID        string                 `json:"id"`
	Type      EventType              `json:"type"`
	Source    string                 `json:"source"`
	Timestamp time.Time              `json:"timestamp"`
	Data      map[string]interface{} `json:"data,omitempty"`
	TraceID   string                 `json:"trace_id,omitempty"`
	UserID    string                 `json:"user_id,omitempty"`
}

// EventBus is a simple in-memory pub/sub for domain events.
// In production, this would be backed by Kafka or another event store.
type EventBus struct {
	subscribers map[EventType][]func(Event)
}

// NewEventBus creates a new event bus.
func NewEventBus() *EventBus {
	return &EventBus{
		subscribers: make(map[EventType][]func(Event)),
	}
}

// Subscribe registers a handler for a specific event type.
func (b *EventBus) Subscribe(eventType EventType, handler func(Event)) {
	b.subscribers[eventType] = append(b.subscribers[eventType], handler)
}

// Publish dispatches an event to all subscribers.
func (b *EventBus) Publish(event Event) {
	if handlers, ok := b.subscribers[event.Type]; ok {
		for _, handler := range handlers {
			go func(h func(Event)) {
				defer func() {
					if r := recover(); r != nil {
						fmt.Printf("event handler panic for %s: %v", event.Type, r)
					}
				}()
				h(event)
			}(handler)
		}
	}
}

// Global event bus instance.
var GlobalBus = NewEventBus()
