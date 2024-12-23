package upload

import "context"

type UploadService interface {
	UploadFile(ctx context.Context, file interface{}, filename string) (string, error)
}
