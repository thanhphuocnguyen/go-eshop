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

func ParsePgNumeric(value float64) (pgtype.Numeric, error) {
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
