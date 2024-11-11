package main

import (
	"github.com/Darkhackit/events/api"
	"github.com/Darkhackit/events/events"
)

func main() {

	events.InitialDispatcher()

	api.Start()
}
