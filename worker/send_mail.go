package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Darkhackit/events/domain"
	"github.com/Darkhackit/events/mail"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rs/zerolog/log"
)

const (
	TaskSendWelcomeMail = "task:send_welcome_mail"
)

type PayloadSendWelcomeEmail struct {
	User domain.User
}

func (distributor *RedisTaskDistributor) DistributeTaskSendWelcome(ctx context.Context, payload *PayloadSendWelcomeEmail, opts ...asynq.Option) error {

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	task := asynq.NewTask(TaskSendWelcomeMail, jsonPayload, opts...)

	enqueueTask, err := distributor.client.EnqueueContext(ctx, task)
	if err != nil {
		return err
	}
	log.Info().
		Str("type", task.Type()).
		Bytes("payload", task.Payload()).
		Str("queue", enqueueTask.Queue).
		Int("max_retry", enqueueTask.MaxRetry).
		Msg("Enqueue task")

	return nil

}

func (processor *RedisTaskProcessor) ProcessTaskSendWelcomeMail(ctx context.Context, task *asynq.Task) error {
	var payload PayloadSendWelcomeEmail
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", asynq.SkipRetry)
	}
	_, err := processor.q.GetUser(ctx, pgtype.Text{
		String: payload.User.Username,
		Valid:  payload.User.Username != "",
	})
	if err != nil {
		return fmt.Errorf("failed to get user from queue: %w", err)
	}

	err = mail.SendWelcomeMail(payload.User.Email, payload.User.Username)
	if err != nil {
		return err
	}
	log.Info().
		Str("type", task.Type()).
		Bytes("payload", task.Payload()).
		Str("email", payload.User.Email).
		Msg("Enqueue task")

	return nil
}
