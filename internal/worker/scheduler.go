package worker

import (
	"time"

	"github.com/hibiken/asynq"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
	app_logger "github.com/thanhphuocnguyen/go-eshop/pkg/logger"
)

type RedisTaskScheduler struct {
	scheduler *asynq.Scheduler
	repo      repository.Store
}

// SendEmailConfirmation implements TaskScheduler.
func (r *RedisTaskScheduler) SendEmailConfirmation() error {
	panic("unimplemented")
}

// SendOrderCreatedEmail implements TaskScheduler.
func (r *RedisTaskScheduler) SendOrderCreatedEmail() error {
	panic("unimplemented")
}

type TaskScheduler interface {
	SendEmailConfirmation() error
	SendOrderCreatedEmail() error
}

func NewRedisTaskScheduler(redisOtp asynq.RedisClientOpt, postgres repository.Store) TaskScheduler {
	// Example of using America/Los_Angeles timezone instead of the default UTC timezone.
	loc, err := time.LoadLocation("Vietnam/Ho_Chi_Minh")
	if err != nil {
		panic(err)
	}
	scheduler := asynq.NewScheduler(
		redisOtp,
		&asynq.SchedulerOpts{
			Location: loc,
			Logger:   app_logger.NewLogger(nil),
		},
	)

	return &RedisTaskScheduler{scheduler, postgres}
}
