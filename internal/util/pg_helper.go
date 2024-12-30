package util

import (
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func GetPgTypeText(value string) pgtype.Text {
	return pgtype.Text{
		String: value,
		Valid:  true,
	}
}

func ParsePgTypeNumber(value float64) (pgtype.Numeric, error) {
	price := pgtype.Numeric{}
	if err := price.Scan(fmt.Sprintf("%.2f", value)); err != nil {
		return price, err
	}
	return price, nil
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

func GetPgTypeInt4(value int32) pgtype.Int4 {
	return pgtype.Int4{
		Int32: value,
		Valid: true,
	}
}

func GetPgTypeInt8(value int64) pgtype.Int8 {
	return pgtype.Int8{
		Int64: value,
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
