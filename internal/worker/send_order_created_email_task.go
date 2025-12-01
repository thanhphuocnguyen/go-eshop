package worker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
	"github.com/thanhphuocnguyen/go-eshop/internal/utils"
)

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
		RawJSON("payload", task.Payload()).
		Str("queue", info.Queue).
		Int("max_retry", info.MaxRetry).
		Msg("sent order created task!!")
	return nil
}

func (p *RedisTaskProcessor) ProcessSendOrderCreatedEmail(ctx context.Context, task *asynq.Task) error {
	var payload PayloadSendOrderCreatedEmailTask
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return err
	}

	payment, err := p.repo.GetPaymentByID(ctx, payload.PaymentID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return fmt.Errorf("could not find payment transaction: %w", asynq.SkipRetry)
		}

		return fmt.Errorf("could not get payment transaction: %w", asynq.SkipRetry)
	}

	order, err := p.repo.GetOrder(ctx, payment.OrderID)
	if err != nil {
		return fmt.Errorf("could not get order details: %w", asynq.SkipRetry)
	}
	orderItems, err := p.repo.GetOrderItems(ctx, payment.OrderID)

	if err != nil {
		return fmt.Errorf("could not get order items: %w", asynq.SkipRetry)
	}

	items := make([]OrderCreatedItems, 0)

	for _, item := range orderItems {
		price, _ := item.LineTotalSnapshot.Float64Value()
		if len(items) == 0 || items[len(items)-1].Name != item.ProductName {
			items = append(items, OrderCreatedItems{
				Name:  item.ProductName,
				Price: price.Float64,
				Qty:   int(item.Quantity),
			})
		}
	}

	user, err := p.repo.GetUserByID(ctx, order.UserID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return fmt.Errorf("could not find user: %w", asynq.SkipRetry)
		}
		return fmt.Errorf("could not get user: %w", asynq.SkipRetry)
	}

	price, _ := payment.Amount.Float64Value()
	emailData := OrderCreatedEmailData{
		OrderID:  order.ID,
		Total:    price.Float64,
		FullName: user.FirstName + " " + user.LastName,
		Items:    items,
	}

	body, err := utils.ParseHtmlTemplate("./static/templates/order-created.html", emailData)

	if err != nil {
		log.Err(err).Msg("could not parse html template")
	}

	err = p.mailer.Send("Order Confirmation - #"+order.ID.String(), body, []string{user.Email}, nil, nil, nil)
	if err != nil {
		return fmt.Errorf("could not send email: %w", err)
	}

	return nil
}
