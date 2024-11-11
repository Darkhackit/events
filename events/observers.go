package events

import (
	"fmt"
	"github.com/Darkhackit/events/domain"
)

type UserCreatedEvent struct {
	User domain.User
}

func (e UserCreatedEvent) Name() string {
	return "UserCreated"
}

func SendWelcomeEmail(event Event) {
	if e, ok := event.(UserCreatedEvent); ok {
		fmt.Printf("Sending welcome email to user with ID: %d\n", e.User.Email)
	}
}

func LogUserCreation(event Event) {
	if e, ok := event.(UserCreatedEvent); ok {
		fmt.Printf("Logging creation of user with ID: %d\n", e.User.Email)
	}
}
