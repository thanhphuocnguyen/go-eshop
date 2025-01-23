package worker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
	"github.com/thanhphuocnguyen/go-eshop/internal/utils"
)

type PayloadSendOrderCreatedEmailTask struct {
	PaymentID string    `json:"payment_id"`
	OrderID   uuid.UUID `json:"order_id"`
}

type OrderCreatedItems struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	Qty   int     `json:"qty"`
}
type OrderCreatedEmailData struct {
	OrderID  uuid.UUID           `json:"order_id"`
	Total    float64             `json:"total"`
	FullName string              `json:"full_name"`
	Items    []OrderCreatedItems `json:"items"`
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
	var payload PayloadSendOrderCreatedEmailTask
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return err
	}
	payment, err := p.repo.GetPaymentTransactionByOrderID(ctx, payload.OrderID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return fmt.Errorf("could not find payment transaction: %w", asynq.SkipRetry)
		}

		return fmt.Errorf("could not get payment transaction: %w", asynq.SkipRetry)
	}
	if payment.PaymentID != payload.PaymentID {
		return fmt.Errorf("payment id mismatch: %w", asynq.SkipRetry)
	}

	orderItems, err := p.repo.GetOrderDetails(ctx, payload.OrderID)
	if err != nil {
		return fmt.Errorf("could not get order details: %w", asynq.SkipRetry)
	}

	items := make([]OrderCreatedItems, 0)

	for _, item := range orderItems {
		price, _ := item.ItemPrice.Float64Value()
		if len(items) == 0 || items[len(items)-1].Name != item.ProductName {
			items = append(items, OrderCreatedItems{
				Name:  item.ProductName,
				Price: price.Float64,
				Qty:   int(item.Quantity),
			})
		}
	}

	user, err := p.repo.GetUserByID(ctx, orderItems[0].UserID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return fmt.Errorf("could not find user: %w", asynq.SkipRetry)
		}
		return fmt.Errorf("could not get user: %w", asynq.SkipRetry)
	}

	price, _ := payment.Amount.Float64Value()
	emailData := OrderCreatedEmailData{
		OrderID:  payload.OrderID,
		Total:    price.Float64,
		FullName: user.Fullname,
		Items:    items,
	}

	body, err := utils.ParseHtmlTemplate("./pkg/mailer/templates/order-created.html", emailData)

	if err != nil {
		log.Err(err).Msg("could not parse html template")
	}

	err = p.mailer.Send("Verify Email", body, []string{user.Email}, nil, nil, nil)
	if err != nil {
		return fmt.Errorf("could not send email: %w", err)
	}

	return nil
}
