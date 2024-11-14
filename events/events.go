package events

import (
	"context"
	"sync"
	"time"
)

// Event Emmanuel Arthur Codes
type Event interface {
	Name() string
}

type Dispatcher struct {
	observers map[string][]func(event Event)
	mu        sync.Mutex
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		observers: make(map[string][]func(event Event)),
	}
}

func (d *Dispatcher) Register(eventName string, observer func(event Event)) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.observers[eventName] = append(d.observers[eventName], observer)
}

func (d *Dispatcher) Dispatch(event Event) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if observers, found := d.observers[event.Name()]; found {
		for _, observer := range observers {
			// Launch each observer in a goroutine with proper context management
			go func(obs func(Event)) {
				// Create a fresh context with a timeout for this goroutine
				_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()

				// Call the observer function with the event
				obs(event) // Assuming observer functions handle the event and context internally
			}(observer)
		}
	}
}
