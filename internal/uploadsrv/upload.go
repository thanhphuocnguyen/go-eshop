package uploadsrv

import "context"

type UploadService interface {
	UploadFile(ctx context.Context, file interface{}, filename string) (string, error)
	RemoveFile(ctx context.Context, publicID string) (string, error)
}
