package events

import "fmt"

var Dispatch *Dispatcher

func InitialDispatcher() {

	Dispatch = NewDispatcher()

	Dispatch.Register("UserCreated", SendWelcomeEmail)
	Dispatch.Register("UserCreated", LogUserCreation)

	fmt.Println("Event Registered")
}
