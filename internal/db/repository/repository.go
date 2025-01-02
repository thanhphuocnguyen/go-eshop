package repository

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/thanhphuocnguyen/go-eshop/config"
)

type Repository interface {
	Querier
	CheckoutCartTx(ctx context.Context, arg CheckoutCartTxParams) (CheckoutCartTxResult, error)
	SetPrimaryAddressTx(ctx context.Context, arg SetPrimaryAddressTxParams) error
	SetPrimaryImageTx(ctx context.Context, arg SetPrimaryImageTxParams) error
	CancelOrderTx(ctx context.Context, params CancelOrderTxParams) error
	Close()
}

const (
	_defaultConnAttempts = 5
	_defaultConnTimeout  = 5 * time.Second
)

type pgRepo struct {
	*Queries
	DbPool       *pgxpool.Pool
	connAttempts int
	connTimeOut  time.Duration
	maxPoolSize  int
}

var once sync.Once
var pg *pgRepo

func GetPostgresInstance(ctx context.Context, cfg config.Config) (Repository, error) {
	var err error
	once.Do(func() {
		pg, err = initializePostgres(ctx, cfg)
	})
	return pg, err
}

func initializePostgres(ctx context.Context, cfg config.Config) (*pgRepo, error) {
	pg = &pgRepo{
		connAttempts: _defaultConnAttempts,
		maxPoolSize:  cfg.MaxPoolSize,
		connTimeOut:  _defaultConnTimeout,
	}

	poolConfig, err := pgxpool.ParseConfig(cfg.DbUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to parse postgres config: %w", err)
	}

	poolConfig.MaxConns = int32(pg.maxPoolSize)

	for pg.connAttempts > 0 {
		pg.DbPool, err = pgxpool.NewWithConfig(ctx, poolConfig)
		if err == nil {
			break
		}
		pg.connAttempts--
		time.Sleep(pg.connTimeOut)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	return pg, nil
}

func (pg *pgRepo) Close() {
	if pg.DbPool != nil {
		pg.DbPool.Close()
	}
}

func (store *pgRepo) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.DbPool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx error: %v, rb error: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit(ctx)
}
