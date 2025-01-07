package worker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
	"github.com/thanhphuocnguyen/go-eshop/internal/utils"
)

type PayloadVerifyEmail struct {
	UserID int64 `json:"user_id"`
}

const (
	VerifyEmailTaskType = "send_verify_email"
)

func (distributor *RedisTaskDistributor) SendVerifyEmail(ctx context.Context, payload *PayloadVerifyEmail, options ...asynq.Option) error {
	marshaled, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("could not marshal payload: %w", err)
	}
	task := asynq.NewTask(VerifyEmailTaskType, marshaled, options...)
	info, err := distributor.client.EnqueueContext(ctx, task)
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

func (processor *RedisTaskProcessor) ProcessSendVerifyEmail(ctx context.Context, t *asynq.Task) error {
	var payload PayloadVerifyEmail
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("could not unmarshal payload: %w", asynq.SkipRetry)
	}

	user, err := processor.repo.GetUserByID(ctx, payload.UserID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return fmt.Errorf("could not find user: %w", asynq.SkipRetry)
		}
		return fmt.Errorf("could not get user: %w", err)
	}

	verifyEmail, err := processor.repo.CreateVerifyEmail(ctx, repository.CreateVerifyEmailParams{
		UserID:     user.UserID,
		Email:      user.Email,
		VerifyCode: utils.RandomString(32),
	})

	if err != nil {
		return fmt.Errorf("could not create verify email: %w", err)
	}

	data, err := os.ReadFile("./pkg/mailer/templates/verify-email.html")

	if err != nil {
		log.Err(err).Msg("could not read email template")
	}

	mailTemplate := string(data)
	verifyLink := fmt.Sprintf("http://%s/v1/verify-email?id=%d&verify_code=%s", processor.cfg.HttpAddr, verifyEmail.ID, verifyEmail.VerifyCode)

	body := strings.Replace(mailTemplate, "[%userName%]", user.Fullname, 1)
	body = strings.Replace(body, "[%verifyLink%]", verifyLink, 1)
	err = processor.mailer.Send("Verify Email", body, []string{user.Email}, nil, nil, nil)
	if err != nil {
		return fmt.Errorf("could not send email: %w", err)
	}

	log.Info().
		Str("username", user.Username).
		Str("email", user.Email).
		Msg("sent verify email to user")
	return nil
}
