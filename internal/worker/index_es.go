package worker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
)

func (d *RedisTaskDistributor) SendIndexProductTask(ctx context.Context, payload *PayloadIndexProduct, options ...asynq.Option) error {
	marshaled, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("could not marshal payload: %w", err)
	}
	task := asynq.NewTask(IndexProductTaskType, marshaled, options...)
	info, err := d.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("could not enqueue task: %w", err)
	}
	log.Info().
		Str("type", task.Type()).
		RawJSON("payload", task.Payload()).
		Str("queue", info.Queue).
		Int("max_retry", info.MaxRetry).
		Msg("sent index product task!!")
	return nil
}

func (d *RedisTaskDistributor) SendUpdateIndexProductTask(ctx context.Context, payload *PayloadIndexProduct, options ...asynq.Option) error {
	marshaled, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("could not marshal payload: %w", err)
	}
	task := asynq.NewTask(UpdateIndexProductTaskType, marshaled, options...)
	info, err := d.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("could not enqueue task: %w", err)
	}
	log.Info().
		Str("type", task.Type()).
		RawJSON("payload", task.Payload()).
		Str("queue", info.Queue).
		Int("max_retry", info.MaxRetry).
		Msg("sent update index product task!!")
	return nil
}

func (d *RedisTaskDistributor) SendDeleteIndexProductTask(ctx context.Context, payload *PayloadIndexProduct, options ...asynq.Option) error {
	marshaled, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("could not marshal payload: %w", err)
	}
	task := asynq.NewTask(DeleteIndexProductTaskType, marshaled, options...)
	info, err := d.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("could not enqueue task: %w", err)
	}
	log.Info().
		Str("type", task.Type()).
		RawJSON("payload", task.Payload()).
		Str("queue", info.Queue).
		Int("max_retry", info.MaxRetry).
		Msg("sent delete index product task!!")
	return nil
}

func (p *RedisTaskProcessor) ProcessIndexProduct(ctx context.Context, task *asynq.Task) error {
	var payload PayloadIndexProduct
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return err
	}

	product, err := p.repo.GetProductForIndexingById(ctx, payload.ProductID)

	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return fmt.Errorf("could not find product: %w", asynq.SkipRetry)
		}

		return fmt.Errorf("could not get product: %w", asynq.SkipRetry)
	}

	err = p.esClient.IndexProduct(repository.GetProductsForIndexingRow{
		ID:               product.ID,
		Name:             product.Name,
		Description:      product.Description,
		BasePrice:        product.BasePrice,
		Categories:       product.Categories,
		ShortDescription: product.ShortDescription,
		Slug:             product.Slug,
		BaseSku:          product.BaseSku,
		IsActive:         product.IsActive,
		ImageUrl:         product.ImageUrl,
		ImageID:          product.ImageID,
		AvgRating:        product.AvgRating,
		Collections:      product.Collections,
		Brand:            product.Brand,
		TotalStock:       product.TotalStock,
		AttributeValues:  product.AttributeValues,
		VariantCount:     product.VariantCount,
		MinPrice:         product.MinPrice,
		UpdatedAt:        product.UpdatedAt,
		CreatedAt:        product.CreatedAt,
	})
	log.Info().Msgf("Indexed product with ID: %s", product.ID.String())
	return nil
}

func (p *RedisTaskProcessor) ProcessUpdateIndexProduct(ctx context.Context, task *asynq.Task) error {
	var payload PayloadIndexProduct
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return err
	}

	product, err := p.repo.GetProductForIndexingById(ctx, payload.ProductID)

	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return fmt.Errorf("could not find product: %w", asynq.SkipRetry)
		}

		return fmt.Errorf("could not get product: %w", asynq.SkipRetry)
	}

	err = p.esClient.UpdateProduct(product.ID.String(), repository.GetProductsForIndexingRow{
		ID:               product.ID,
		Name:             product.Name,
		Description:      product.Description,
		BasePrice:        product.BasePrice,
		Categories:       product.Categories,
		ShortDescription: product.ShortDescription,
		Slug:             product.Slug,
		BaseSku:          product.BaseSku,
		IsActive:         product.IsActive,
		ImageUrl:         product.ImageUrl,
		ImageID:          product.ImageID,
		AvgRating:        product.AvgRating,
		Collections:      product.Collections,
		Brand:            product.Brand,
		TotalStock:       product.TotalStock,
		AttributeValues:  product.AttributeValues,
		VariantCount:     product.VariantCount,
		MinPrice:         product.MinPrice,
		UpdatedAt:        product.UpdatedAt,
		CreatedAt:        product.CreatedAt,
	})
	log.Info().Msgf("Updated indexed product with ID: %s", product.ID.String())
	return nil
}

func (p *RedisTaskProcessor) ProcessDeleteIndexProduct(ctx context.Context, task *asynq.Task) error {
	var payload PayloadIndexProduct
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return err
	}

	product, err := p.repo.GetProductForIndexingById(ctx, payload.ProductID)

	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return fmt.Errorf("could not find payment transaction: %w", asynq.SkipRetry)
		}

		return fmt.Errorf("could not get product: %w", asynq.SkipRetry)
	}

	err = p.esClient.DeleteProduct(product.ID.String())
	log.Info().Msgf("Deleted indexed product with ID: %s", product.ID.String())
	return err
}
