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
	"github.com/thanhphuocnguyen/go-eshop/pkg/elasticsearch"
)

func (d *RedisTaskDistributor) SendIndexProductTask(ctx context.Context, payload PayloadIndexProduct, options ...asynq.Option) error {
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

func (d *RedisTaskDistributor) SendUpdateIndexProductTask(ctx context.Context, payload PayloadUpdateProductIndex, options ...asynq.Option) error {
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

func (d *RedisTaskDistributor) SendDeleteIndexProductTask(ctx context.Context, payload PayloadIndexProduct, options ...asynq.Option) error {
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

	product, err := p.store.GetProductForIndexingById(ctx, payload.ProductID)

	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return fmt.Errorf("could not find product: %w", asynq.SkipRetry)
		}

		return fmt.Errorf("could not get product: %w", asynq.SkipRetry)
	}

	err = p.esClient.IndexProduct(ctx, repository.GetProductsForIndexingRow{
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
	var payload PayloadUpdateProductIndex
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return err
	}
	data := payload.Data

	updateData := elasticsearch.ProductIndexDocument{
		ID: payload.ProductID.String(),
	}

	if data.Name != nil {
		updateData.Name = *data.Name
	}
	if data.Description != nil {
		updateData.Description = *data.Description
	}
	if data.ShortDescription != nil {
		updateData.ShortDescription = data.ShortDescription
	}
	if data.IsActive != nil {
		updateData.IsActive = *data.IsActive
	}
	if data.Slug != nil {
		updateData.Slug = *data.Slug
	}
	if data.BaseSku != nil {
		updateData.Sku = *data.BaseSku
	}
	if data.CategoryIDs != nil {
		ids := make([]uuid.UUID, 0)
		for _, idStr := range *data.CategoryIDs {
			id, err := uuid.Parse(idStr)
			if err != nil {
				return fmt.Errorf("could not parse category id: %w", err)
			}
			ids = append(ids, id)
		}
		categories, err := p.store.GetCategoriesByIDs(ctx, ids)
		names := make([]string, 0)
		for _, category := range categories {
			names = append(names, category.Name)
		}
		if err != nil {
			return fmt.Errorf("could not get categories by ids: %w", err)
		}
		updateData.Categories = names
	}

	if data.CollectionIDs != nil {
		ids := make([]uuid.UUID, 0)
		for _, idStr := range *data.CollectionIDs {
			id, err := uuid.Parse(idStr)
			if err != nil {
				return fmt.Errorf("could not parse collection id: %w", err)
			}
			ids = append(ids, id)
		}
		collections, err := p.store.GetCollectionsByIDs(ctx, ids)
		names := make([]string, 0)
		for _, collection := range collections {
			names = append(names, collection.Name)
		}
		if err != nil {
			return fmt.Errorf("could not get collections by ids: %w", err)
		}
		updateData.Collections = names
	}
	if data.BrandID != nil {
		brand, err := p.store.GetBrandByID(ctx, uuid.MustParse(*data.BrandID))
		if err != nil {
			return fmt.Errorf("could not get brand by id: %w", err)
		}
		updateData.Brand = brand.Name
	}

	err := p.esClient.UpdateProduct(ctx, payload.ProductID.String(), updateData)
	if err != nil {
		return err
	}
	log.Info().Msgf("Updated indexed product with ID: %s", payload.ProductID.String())
	return nil
}

func (p *RedisTaskProcessor) ProcessDeleteIndexProduct(ctx context.Context, task *asynq.Task) error {
	var payload PayloadIndexProduct
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return err
	}

	product, err := p.store.GetProductForIndexingById(ctx, payload.ProductID)

	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return fmt.Errorf("could not find payment transaction: %w", asynq.SkipRetry)
		}

		return fmt.Errorf("could not get product: %w", asynq.SkipRetry)
	}

	err = p.esClient.DeleteProduct(ctx, product.ID.String())
	log.Info().Msgf("Deleted indexed product with ID: %s", product.ID.String())
	return err
}
