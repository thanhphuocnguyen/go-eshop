package utils

import (
	"fmt"
	"math"
	"strings"
	"time"
)

func GetImageName(originFileName string) string {
	return fmt.Sprintf("%d-%s", time.Now().UnixNano(), strings.ReplaceAll(originFileName, " ", "-"))
}

func StandardizeDecimal(num float64) float64 {
	return math.Floor(num*MUL) / EXP
}

func BoolPtr(value bool) *bool {
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

func GetAvgRating(ratingCnt int32, oneStarCnt,
	twoStarCnt, threeStarCnt, fourStarCnt, fiveStarCnt int32) float64 {
	if ratingCnt == 0 {
		return 0
	}
	avg := float64(oneStarCnt*1+twoStarCnt*2+threeStarCnt*3+
		fourStarCnt*4+fiveStarCnt*5) / float64(ratingCnt)
	return avg
}
