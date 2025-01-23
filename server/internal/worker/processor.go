package worker

import (
	"context"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
	"github.com/thanhphuocnguyen/go-eshop/config"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
	logger "github.com/thanhphuocnguyen/go-eshop/pkg/log"
	"github.com/thanhphuocnguyen/go-eshop/pkg/mailer"
)

const (
	QueueCritical = "critical"
	QueueDefault  = "default"
	QueueLow      = "low"
)

type TaskProcessor interface {
	ProcessSendOrderCreatedEmail(ctx context.Context, task *asynq.Task) error
	Start() error
	Shutdown()
}
type RedisTaskProcessor struct {
	asynqServer *asynq.Server
	repo        repository.Repository
	mailer      mailer.EmailSender
	cfg         config.Config
}

func NewRedisTaskProcessor(
	redisOtp asynq.RedisClientOpt,
	postgres repository.Repository,
	mailer mailer.EmailSender,
	cfg config.Config,
) TaskProcessor {
	logger := logger.NewLogger(nil)
	server := asynq.NewServer(redisOtp, asynq.Config{
		Concurrency: 10,
		Queues: map[string]int{
			QueueCritical: 4,
			QueueDefault:  5,
		},
		Logger: logger,
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
	return &RedisTaskProcessor{server, postgres, mailer, cfg}
}

func (p *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()
	// register task handlers
	mux.HandleFunc(OrderCreatedEmailTaskType, p.ProcessSendOrderCreatedEmail)
	mux.HandleFunc(VerifyEmailTaskType, p.ProcessSendVerifyEmail)

	return p.asynqServer.Start(mux)
}

func (p *RedisTaskProcessor) Shutdown() {
	p.asynqServer.Shutdown()
}
