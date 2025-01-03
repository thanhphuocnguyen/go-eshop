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
	CreateOrderTx(ctx context.Context, arg CreateOrderTxParams) (int64, error)
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
var repoInstance *pgRepo

func GetPostgresInstance(ctx context.Context, cfg config.Config) (Repository, error) {
	var err error
	once.Do(func() {
		repoInstance, err = initializePostgres(ctx, cfg)
	})
	return repoInstance, err
}

func initializePostgres(ctx context.Context, cfg config.Config) (*pgRepo, error) {
	repoInstance = &pgRepo{
		connAttempts: _defaultConnAttempts,
		maxPoolSize:  cfg.MaxPoolSize,
		connTimeOut:  _defaultConnTimeout,
	}

	poolConfig, err := pgxpool.ParseConfig(cfg.DbUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to parse postgres config: %w", err)
	}

	poolConfig.MaxConns = int32(repoInstance.maxPoolSize)

	for repoInstance.connAttempts > 0 {
		repoInstance.DbPool, err = pgxpool.NewWithConfig(ctx, poolConfig)
		if err == nil {
			break
		}
		repoInstance.connAttempts--
		time.Sleep(repoInstance.connTimeOut)
	}
	err = repoInstance.DbPool.Ping(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}
	repoInstance.Queries = New(repoInstance.DbPool)

	return repoInstance, nil
}

func (repoConn *pgRepo) Close() {
	if repoConn.DbPool != nil {
		repoConn.DbPool.Close()
		repoConn.DbPool = nil
	}
	repoInstance = nil
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
