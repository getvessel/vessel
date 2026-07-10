package events

import "sync"

type Event string

const (
	ProjectDeployed Event = "project.deployed"
	BackupFailed    Event = "backup.failed"
	MemberInvited   Event = "member.invited"
)

type Listener func(payload interface{})

type EventBus struct {
	listeners map[Event][]Listener
	mu        sync.RWMutex
}

func NewEventBus() *EventBus {
	return &EventBus{
		listeners: make(map[Event][]Listener),
	}
}

func (b *EventBus) Subscribe(event Event, listener Listener) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.listeners[event] = append(b.listeners[event], listener)
}

func (b *EventBus) Publish(event Event, payload interface{}) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	for _, listener := range b.listeners[event] {
		// Run each listener asynchronously so it doesn't block the caller
		go listener(payload)
	}
}
