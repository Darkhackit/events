package events

import (
	"context"
	"fmt"
	"github.com/Darkhackit/events/domain"
	"github.com/Darkhackit/events/worker"
	"github.com/hibiken/asynq"
	"time"
)

type UserCreatedEvent struct {
	User            domain.User
	TaskDistributor worker.TaskDistributor
}

func (e UserCreatedEvent) Name() string {
	return "UserCreated"
}

func SendWelcomeEmail(event Event) {
	if e, ok := event.(UserCreatedEvent); ok {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		taskPayload := &worker.PayloadSendWelcomeEmail{User: e.User}
		ops := []asynq.Option{
			asynq.MaxRetry(5),
			asynq.ProcessIn(10 * time.Second),
			asynq.Queue(worker.QueueCritical),
		}
		err := e.TaskDistributor.DistributeTaskSendWelcome(ctx, taskPayload, ops...)
		if err != nil {
			panic(err)
		}

	}
}

func LogUserCreation(event Event) {
	if e, ok := event.(UserCreatedEvent); ok {
		fmt.Printf("Logging creation of user with ID: %d\n", e.User.Email)
	}
}
