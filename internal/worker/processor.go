package worker

import (
	"context"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
)

const (
	QueueCritical = "critical"
	QueueDefault  = "default"
)

type TaskProcessor interface {
	ProcessSendOrderCreatedEmail(ctx context.Context, task *asynq.Task) error
	Start() error
	Shutdown()
}
type RedisTaskProcessor struct {
	server *asynq.Server
	repo   repository.Repository
}

func NewRedisTaskProcessor(redisOtp asynq.RedisClientOpt, postgres repository.Repository) TaskProcessor {
	server := asynq.NewServer(redisOtp, asynq.Config{
		Concurrency: 10,
		Queues: map[string]int{
			QueueCritical: 4,
			QueueDefault:  5,
		},
		HealthCheckFunc: func(err error) {
			if err != nil {
				log.Error().Err(err).Msg("error with task processor")
			}
		},
		ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
			log.Error().
				Err(err).
				Str("task_type", task.Type()).
				Bytes("payload", task.Payload()).
				Msg("error processing task")
		}),
	})
	return &RedisTaskProcessor{server, postgres}
}

func (p *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()
	// register task handlers
	mux.HandleFunc(OrderCreatedEmailTaskType, p.ProcessSendOrderCreatedEmail)

	return p.server.Start(mux)
}

func (p *RedisTaskProcessor) Shutdown() {
	p.server.Shutdown()
}
