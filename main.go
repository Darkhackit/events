package main

import (
	"fmt"
	"github.com/Darkhackit/events/events"
	"time"
)

type SendWelcomeEmail struct{}
type LogUserCreation struct{}
type UpdateBalance struct {
}

func (s SendWelcomeEmail) Handle(event events.Event) {
	if e, ok := event.(events.UserCreatedEvent); ok {
		fmt.Printf("Sendind welcome email to yser with ID %d\n", e.UserID)
	}
}

func (u UpdateBalance) Handle(event events.Event) {
	if e, ok := event.(events.UserBalanceUpdatedEvent); ok {
		fmt.Printf("Sendind update account balance to yser with ID %d\n", e.UserID)
	}
}

func (s LogUserCreation) Handle(event events.Event) {
	if e, ok := event.(events.UserCreatedEvent); ok {
		fmt.Printf("Sendind welcome email to yser with ID %d\n", e.UserID)
	}
}

func main() {
	dispatcher := events.NewDispatcher()
	dispatcher.Register("UserCreated", SendWelcomeEmail{})
	dispatcher.Register("UserCreated", LogUserCreation{})
	dispatcher.Register("UserBalanceUpdated", UpdateBalance{})

	userEvent := events.UserCreatedEvent{UserID: 54}
	ev := events.UserBalanceUpdatedEvent{
		UserID: 100,
		Amount: 500,
	}

	dispatcher.Dispatch(ev)

	dispatcher.Dispatch(userEvent)

	time.Sleep(1 * time.Second)
}
