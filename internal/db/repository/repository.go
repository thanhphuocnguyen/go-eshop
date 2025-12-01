package repository

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	"github.com/thanhphuocnguyen/go-eshop/config"
)

type Repository interface {
	Querier
	CheckoutCartTx(ctx context.Context, arg CheckoutCartTxArgs) (CreatePaymentResult, error)
	SetPrimaryAddressTx(ctx context.Context, arg SetPrimaryAddressTxArgs) error
	CancelOrderTx(ctx context.Context, params CancelOrderTxArgs) (uuid.UUID, error)
	RefundOrderTx(ctx context.Context, params RefundOrderTxArgs) error
	VerifyEmailTx(ctx context.Context, arg VerifyEmailTxArgs) error
	CreateProductTx(ctx context.Context, arg CreateProductTxArgs) (Product, error)
	UpdateProductTx(ctx context.Context, arg UpdateProductTxArgs) (Product, error)
	QueryRaw(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error)
	VoteHelpfulRatingTx(ctx context.Context, arg VoteHelpfulRatingTxArgs) (uuid.UUID, error)
	UpdateDiscountTx(ctx context.Context, id uuid.UUID, arg UpdateDiscountTxArgs) error
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
var mu sync.Mutex

func GetPostgresInstance(ctx context.Context, cfg config.Config) (Repository, error) {
	mu.Lock()
	defer mu.Unlock()

	// If instance is nil, reset once to allow re-initialization
	if repoInstance == nil {
		once = sync.Once{}
	}

	var err error
	once.Do(func() {
		repoInstance, err = initializePostgres(ctx, cfg)
	})

	if err != nil {
		return nil, fmt.Errorf("failed to initialize postgres: %w", err)
	}

	if repoInstance.DbPool == nil {
		return nil, fmt.Errorf("database pool is not initialized")
	}

	err = repoInstance.DbPool.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	log.Info().Msg("Connected to postgres")
	return repoInstance, nil
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
	// Initialize queries
	repoInstance.Queries = New(repoInstance.DbPool)
	return repoInstance, nil
}

func (repoConn *pgRepo) Close() {
	mu.Lock()
	defer mu.Unlock()

	if repoConn != nil && repoConn.DbPool != nil {
		// Create a context with timeout for graceful shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Close the pool in a goroutine to avoid blocking
		done := make(chan struct{})
		go func() {
			defer close(done)
			repoConn.DbPool.Close()
		}()

		// Wait for close to complete or timeout
		select {
		case <-done:
			log.Info().Msg("Database pool closed successfully")
		case <-ctx.Done():
			log.Warn().Msg("Database pool close timed out")
		}

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
func (store *pgRepo) QueryRaw(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error) {
	return store.DbPool.Query(ctx, query, args...)
}
