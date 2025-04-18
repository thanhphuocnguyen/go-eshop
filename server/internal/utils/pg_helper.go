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

func GetPgTypeText(value string) pgtype.Text {
	return pgtype.Text{
		String: value,
		Valid:  true,
	}
}
func GetPgTypeBool(value bool) pgtype.Bool {
	return pgtype.Bool{
		Bool:  value,
		Valid: true,
	}
}

func GetPgTypeTimestamp(value time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{
		Time:  value,
		Valid: true,
	}
}

func GetPgTypeInt2(value int16) pgtype.Int2 {
	return pgtype.Int2{
		Int16: value,
		Valid: true,
	}
}
func GetPgTypeInt4(value int32) pgtype.Int4 {
	return pgtype.Int4{
		Int32: value,
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

func GetPgTypeInt8(value int64) pgtype.Int8 {
	return pgtype.Int8{
		Int64: value,
		Valid: true,
	}
}

func GetPgTypeUUID(value [16]byte) pgtype.UUID {
	return pgtype.UUID{
		Bytes: value,
		Valid: true,
	}
}

func GetPgTypeFloat4(value float32) pgtype.Float4 {
	return pgtype.Float4{
		Float32: value,
		Valid:   true,
	}
}

func GetPgTypeFloat8(value float64) pgtype.Float8 {
	return pgtype.Float8{
		Float64: value,
		Valid:   true,
	}
}

func GetPgNumericFromFloat(value float64) pgtype.Numeric {
	return pgtype.Numeric{
		Int:   big.NewInt(int64(value * MUL)),
		Exp:   -EXP,
		Valid: true,
	}
}
