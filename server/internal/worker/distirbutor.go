package worker

import (
	"context"

	"github.com/hibiken/asynq"
)

type TaskDistributor interface {
	SendOrderCreatedEmailTask(ctx context.Context, payload *PayloadSendOrderCreatedEmailTask, options ...asynq.Option) error
	SendVerifyEmail(ctx context.Context, payload *PayloadVerifyEmail, options ...asynq.Option) error
	Shutdown() error
}

type RedisTaskDistributor struct {
	client *asynq.Client
}

func NewRedisTaskDistributor(redisOpt asynq.RedisClientOpt) TaskDistributor {
	client := asynq.NewClient(redisOpt)
	return &RedisTaskDistributor{client}
}

func (d *RedisTaskDistributor) Shutdown() error {
	return d.client.Close()
}
