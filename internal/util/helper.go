package util

import (
	"fmt"

	"github.com/google/uuid"
)

func GetImageName(originFileName string, userID, productID int64) string {
	guid := uuid.New()
	return fmt.Sprintf("user_%d_product_%d_%s_%s", userID, productID, guid.String(), originFileName)
}

func GetImageURL(fileName string) string {
	return fmt.Sprintf("assets/images/%s", fileName)
}
