package events

import (
	"sync"
)

type UserCreatedEvent struct {
	UserID int
}
type UserBalanceUpdatedEvent struct {
	UserID int
	Amount float64
}

type Event interface {
	Name() string
}
type Observer interface {
	Handle(event Event)
}

func (e UserCreatedEvent) Name() string {
	return "UserCreated"
}
func (e UserBalanceUpdatedEvent) Name() string {
	return "UserBalanceUpdated"
}

type Dispatcher struct {
	observers map[string][]Observer
	mu        sync.Mutex
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		observers: make(map[string][]Observer),
	}
}

func (d *Dispatcher) Register(eventName string, observer Observer) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.observers[eventName] = append(d.observers[eventName], observer)
}

func (d *Dispatcher) Dispatch(event Event) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if observers, found := d.observers[event.Name()]; found {
		for _, observer := range observers {
			go observer.Handle(event)
		}
	}
}
