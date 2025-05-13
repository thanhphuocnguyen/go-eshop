package api

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/thanhphuocnguyen/go-eshop/internal/utils"
)

func GetImageName(originFileName string) string {
	return fmt.Sprintf("%d-%s", time.Now().UnixNano(), strings.ReplaceAll(originFileName, " ", "-"))
}

func StandardizeDecimal(num float64) float64 {
	return math.Floor(num*utils.MUL) / utils.EXP
}

func String(value string) *string {
	return &value
}

func Bool(value bool) *bool {
	return &value
}
func StringPtr(value string) *string {
	return &value
}
func Int32Ptr(value int32) *int32 {
	return &value
}
func Int64Ptr(value int64) *int64 {
	return &value
}

func CalculateTotalPages(total int64, pageSize int64) int64 {
	if total == 0 {
		return 0
	}
	if total%pageSize == 0 {
		return total / pageSize
	}
	return total/pageSize + 1
}

func TimeDurationPtr(value time.Duration) *time.Duration {
	return &value
}
