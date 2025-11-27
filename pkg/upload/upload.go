package upload

import "context"

type CdnUploader interface {
	Upload(ctx context.Context, file interface{}) (publicID string, url string, err error)
	Remove(ctx context.Context, ID string) (message string, err error)
}
