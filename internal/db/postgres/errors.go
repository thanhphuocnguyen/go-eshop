package postgres

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

var ErrorRecordNotFound = pgx.ErrNoRows
var ErrorForeignKeyViolation = &pgconn.PgError{Code: ForeignKeyViolation}
var ErrorUniqueViolation = &pgconn.PgError{Code: UniqueViolation}
var ErrorTimeout = &pgconn.PgError{Code: Timeout}
var ErrorDeadlockDetected = &pgconn.PgError{Code: DeadLock}

func ErrorCode(err error) string {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code
	}
	return ""
}
