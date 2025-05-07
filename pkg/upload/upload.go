package upload

import "context"

type UploadService interface {
	UploadFile(ctx context.Context, file interface{}) (publicID string, url string, err error)
	RemoveFile(ctx context.Context, ID string) (message string, err error)
}
