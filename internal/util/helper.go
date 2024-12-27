package util

import (
	"fmt"
	"strings"
	"time"
)

func GetImageName(originFileName string) string {
	return fmt.Sprintf("%d-%s", time.Now().UnixNano(), strings.ReplaceAll(originFileName, " ", "-"))
}

func GetImageURL(fileName string) string {
	return fmt.Sprintf("assets/images/%s", fileName)
}
