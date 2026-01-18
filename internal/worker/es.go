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

func (p *RedisTaskProcessor) ProcessIndexProduct(ctx context.Context, task *asynq.Task) error {
	var payload PayloadIndexProduct
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return err
	}

	product, err := p.repo.GetProductDetail(ctx, repository.GetProductDetailParams{
		ID: payload.ProductID,
	})

	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return fmt.Errorf("could not find payment transaction: %w", asynq.SkipRetry)
		}

		return fmt.Errorf("could not get payment transaction: %w", asynq.SkipRetry)
	}

	categories := make([]string, 0)
	parsedCats := make([]repository.Category, 0)
	err = json.Unmarshal(product.Categories, &parsedCats)
	for _, category := range parsedCats {
		categories = append(categories, category.Name)
	}
	parsedCollections := make([]repository.Collection, 0)
	err = json.Unmarshal(product.Collections, &parsedCollections)
	collections := make([]string, 0)
	for _, collection := range parsedCollections {
		collections = append(collections, collection.Name)
	}
	brand := repository.Brand{}
	err = json.Unmarshal(product.Brand, &brand)
	if err != nil {
		return fmt.Errorf("could not unmarshal brand: %w", asynq.SkipRetry)
	}
	err = p.esClient.IndexProduct(repository.GetProductsForIndexingRow{
		ID:               product.ID,
		Name:             product.Name,
		Description:      product.Description,
		BasePrice:        product.BasePrice,
		Categories:       categories,
		ShortDescription: product.ShortDescription,
		Slug:             product.Slug,
		BaseSku:          product.BaseSku,
		IsActive:         product.IsActive,
		ImageUrl:         product.ImageUrl,
		ImageID:          product.ImageID,
		AvgRating:        product.AvgRating,
		Collections:      collections,
		Brand:            brand.Name,
		UpdatedAt:        product.UpdatedAt,
		CreatedAt:        product.CreatedAt,
	})
	return nil
}
