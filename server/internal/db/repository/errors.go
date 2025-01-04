package repository

import (
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	ForeignKeyViolation = "23503"
	UniqueViolation     = "23505"
	Timeout             = "57014"
	DeadLock            = "40P01"
)

var ErrRecordNotFound = pgx.ErrNoRows
var ErrForeignKeyViolation = &pgconn.PgError{Code: ForeignKeyViolation}
var ErrUniqueViolation = &pgconn.PgError{Code: UniqueViolation}
var ErrTimeout = &pgconn.PgError{Code: Timeout}
var ErrDeadlockDetected = &pgconn.PgError{Code: DeadLock}

func ErrorCode(err error) string {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code
	}
	return ""
}
