package utils

import (
	"errors"
	"fmt"
	"math"
	"reflect"
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

// Slugify converts a string to a URL-friendly slug
func Slugify(s string) string {
	// Convert to lowercase
	s = strings.ToLower(s)
	// Replace spaces and special characters with hyphens
	s = strings.ReplaceAll(s, " ", "-")
	s = strings.ReplaceAll(s, "&", "and")
	// Remove other non-alphanumeric characters except hyphens
	var result strings.Builder
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			result.WriteRune(r)
		}
	}
	s = result.String()
	// Remove multiple consecutive hyphens
	for strings.Contains(s, "--") {
		s = strings.ReplaceAll(s, "--", "-")
	}
	// Trim hyphens from start and end
	return strings.Trim(s, "-")
}

func SetField(obj interface{}, name string, value interface{}) error {
	structValue := reflect.ValueOf(obj).Elem()
	structFieldValue := structValue.FieldByName(name)

	if !structFieldValue.IsValid() {
		return fmt.Errorf("No such field: %s in obj", name)
	}

	if !structFieldValue.CanSet() {
		return fmt.Errorf("Cannot set %s field value", name)
	}

	structFieldType := structFieldValue.Type()
	val := reflect.ValueOf(value)
	if structFieldType != val.Type() {
		return errors.New("Provided value type didn't match obj field type")
	}

	structFieldValue.Set(val)
	return nil
}

type MyStruct struct {
	Name string
	Age  int64
}
