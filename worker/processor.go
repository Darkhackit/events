package worker

import (
	"context"
	db "github.com/Darkhackit/events/db/sqlc"
	"github.com/hibiken/asynq"
)

const (
	QueueCritical = "critical"
	QueueDefault  = "default"
	QueueLow      = "low"
)

type TaskProcessor interface {
	Start() error
	ProcessTaskSendWelcomeMail(ctx context.Context, task *asynq.Task) error
	Stop()
}

type RedisTaskProcessor struct {
	server *asynq.Server
	q      db.Queries
}

func (processor *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()
	mux.HandleFunc(TaskSendWelcomeMail, processor.ProcessTaskSendWelcomeMail)

	err := processor.server.Start(mux)
	if err != nil {
		return err
	}
	return nil
}
func (processor *RedisTaskProcessor) Stop() {

}

func NewRedisTaskProcessor(redisOpt asynq.RedisClientOpt, q db.Queries) TaskProcessor {

	server := asynq.NewServer(redisOpt, asynq.Config{
		Queues: map[string]int{
			QueueCritical: 8,
			QueueDefault:  4,
			QueueLow:      3,
		},
	})

	return &RedisTaskProcessor{server: server, q: q}

}
