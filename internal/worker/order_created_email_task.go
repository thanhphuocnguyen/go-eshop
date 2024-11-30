package worker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

type PayloadSendOrderCreatedEmailTask struct {
	OrderID int
	// Products
	Total    float64
	FullName string
	Email    string
	Phone    string
}

func (d *RedisTaskDistributor) SendOrderCreatedEmailTask(ctx context.Context, payload *PayloadSendOrderCreatedEmailTask, options ...asynq.Option) error {
	marshaled, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("could not marshal payload: %w", err)
	}

	task := asynq.NewTask(OrderCreatedEmailTaskType, marshaled, options...)
	info, err := d.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("could not enqueue task: %w", err)
	}
	log.Info().
		Str("type", task.Type()).
		Bytes("payload", task.Payload()).
		Str("queue", info.Queue).
		Int("max_retry", info.MaxRetry).
		Msg("task enqueued")

	return nil
}

func (p *RedisTaskProcessor) ProcessSendOrderCreatedEmail(ctx context.Context, task *asynq.Task) error {
	var payload struct {
		OrderID int
	}
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return err
	}
	return nil
}
