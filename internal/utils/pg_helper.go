package utils

import (
	"math/big"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const (
	MUL = 100
	EXP = 2
)

func GetPgTypeTimestamp(value time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{
		Time:  value,
		Valid: true,
	}
}

func GetPgTypeUUIDFromString(value string) pgtype.UUID {
	uuid, err := uuid.Parse(value)
	if err != nil {
		return pgtype.UUID{
			Valid: false,
		}
	}
	return pgtype.UUID{
		Bytes: uuid,
		Valid: true,
	}
}

func GetPgTypeUUID(value [16]byte) pgtype.UUID {
	return pgtype.UUID{
		Bytes: value,
		Valid: true,
	}
}

func GetPgNumericFromFloat(value float64) pgtype.Numeric {
	return pgtype.Numeric{
		Int:   big.NewInt(int64(value * MUL)),
		Exp:   -EXP,
		Valid: true,
	}
}

func StandardizeDecimal(num float64) float64 {
	return num * MUL / EXP
}
