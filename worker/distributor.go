package worker

import "context"
import "github.com/hibiken/asynq"

type TaskDistributor interface {
	DistributeTaskSendWelcome(
		ctx context.Context,
		payload *PayloadSendWelcomeEmail,
		opts ...asynq.Option,
	) error
}

type RedisTaskDistributor struct {
	client *asynq.Client
}

func NewRedisTaskDistributor(redisOpt asynq.RedisClientOpt) TaskDistributor {

	client := asynq.NewClient(redisOpt)

	return &RedisTaskDistributor{client: client}

}
