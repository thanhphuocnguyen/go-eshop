package upload

import "context"

type UploadService interface {
	UploadFile(ctx context.Context, file interface{}, ID string) (string, error)
	RemoveFile(ctx context.Context, ID string) (string, error)
}
