package util

import (
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

func GetPgTypeText(value string) pgtype.Text {
	return pgtype.Text{
		String: value,
		Valid:  true,
	}
}

func ParsePgNumeric(value float32) (pgtype.Numeric, error) {

	price := pgtype.Numeric{}
	if err := price.Scan(fmt.Sprintf("%.2f", value)); err != nil {
		return price, err
	}
	return price, nil
}
