package postgres

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/thanhphuocnguyen/go-eshop/config"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/sqlc"
)

type Store interface {
	sqlc.Querier
	CheckoutCartTx(ctx context.Context, arg CheckoutCartParams) (CheckoutCartTxResult, error)
	Close()
}

const (
	_defaultConnAttempts = 5
	_defaultConnTimeout  = 5 * time.Second
)

type Postgres struct {
	connAttempts int
	connTimeOut  time.Duration
	maxPoolSize  int
	*sqlc.Queries
	Pool *pgxpool.Pool
}

var once sync.Once
var pg *Postgres

func GetPostgresInstance(ctx context.Context, cfg config.Config) (Store, error) {
	var err error
	once.Do(func() {
		pg, err = initializePostgres(ctx, cfg)
	})
	return pg, err
}

func initializePostgres(ctx context.Context, cfg config.Config) (*Postgres, error) {
	pg = &Postgres{
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
		pg.Pool, err = pgxpool.NewWithConfig(ctx, poolConfig)
		if err == nil {
			break
		}
		pg.connAttempts--
		time.Sleep(pg.connTimeOut)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	pg.Queries = sqlc.New(pg.Pool)
	return pg, nil
}

func (pg *Postgres) Close() {
	if pg.Pool != nil {
		pg.Pool.Close()
	}
}

func (store *Postgres) execTx(ctx context.Context, fn func(*sqlc.Queries) error) error {
	tx, err := store.Pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	q := sqlc.New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx error: %v, rb error: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit(ctx)
}
