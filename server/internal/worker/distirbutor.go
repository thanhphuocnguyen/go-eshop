package worker

import (
	"context"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

type TaskDistributor interface {
	SendOrderCreatedEmailTask(ctx context.Context, payload *PayloadSendOrderCreatedEmailTask, options ...asynq.Option) error
	SendVerifyAccountEmail(ctx context.Context, payload *PayloadVerifyEmail, options ...asynq.Option) error
	Shutdown() error
}

type RedisTaskDistributor struct {
	client *asynq.Client
}

func NewRedisTaskDistributor(redisOpt asynq.RedisClientOpt) TaskDistributor {
	client := asynq.NewClient(redisOpt)
	err := client.Ping()
	if err != nil {
		panic(err)
	}

	log.Info().Msg("Distributor connected to Redis!!!")
	return &RedisTaskDistributor{client}
}

func (d *RedisTaskDistributor) Shutdown() error {
	return d.client.Close()
}
