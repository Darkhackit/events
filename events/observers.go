package events

import (
	"fmt"
	"github.com/Darkhackit/events/domain"
	"github.com/Darkhackit/events/mail"
)

type UserCreatedEvent struct {
	User domain.User
}

func (e UserCreatedEvent) Name() string {
	return "UserCreated"
}

func SendWelcomeEmail(event Event) {
	if e, ok := event.(UserCreatedEvent); ok {
		err := mail.SendWelcomeMail(e.User.Email, e.User.Username)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func LogUserCreation(event Event) {
	if e, ok := event.(UserCreatedEvent); ok {
		fmt.Printf("Logging creation of user with ID: %d\n", e.User.Email)
	}
}
