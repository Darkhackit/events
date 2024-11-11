package events

import (
	"sync"
)

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
			go observer(event)
		}
	}
}
