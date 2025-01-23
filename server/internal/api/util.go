package api

import (
	"fmt"
	"strings"
	"time"

	"github.com/thanhphuocnguyen/go-eshop/internal/utils"
)

func GetImageName(originFileName string) string {
	return fmt.Sprintf("%d-%s", time.Now().UnixNano(), strings.ReplaceAll(originFileName, " ", "-"))
}

func StandardizeDecimal(num float64) float64 {
	return float64(int64(num*utils.MUL)) / utils.EXP
}
