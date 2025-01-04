package upload

import (
	"context"
	"fmt"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/cloudinary/cloudinary-go/v2/logger"
	"github.com/thanhphuocnguyen/go-eshop/config"
	"github.com/thanhphuocnguyen/go-eshop/pkg/log"
)

type CloudinaryUploadService struct {
	cld *cloudinary.Cloudinary
	cfg config.Config
}

func NewCloudinaryUploadService(cfg config.Config) UploadService {
	cld, err := cloudinary.NewFromURL(cfg.CloudinaryUrl)
	cld.Logger = &logger.Logger{
		Writer: log.NewLogger(nil),
	}

	if err != nil {
		panic(err)
	}

	return &CloudinaryUploadService{
		cld: cld,
		cfg: cfg,
	}
}

func (s *CloudinaryUploadService) UploadFile(ctx context.Context, file interface{}, filename string) (string, error) {
	uploadResult, err := s.cld.Upload.Upload(ctx, file, uploader.UploadParams{
		PublicID: filename,
		Folder:   s.cfg.CloudinaryFolder,
	})
	if err != nil {
		return "", err
	}
	return uploadResult.SecureURL, nil
}

func (s *CloudinaryUploadService) RemoveFile(ctx context.Context, publicID string) (string, error) {
	publicIDWithPath := fmt.Sprintf("%s/%s", s.cfg.CloudinaryFolder, publicID)
	invalidate := true
	result, err := s.cld.Upload.Destroy(ctx, uploader.DestroyParams{PublicID: publicIDWithPath, Invalidate: &invalidate})
	if err != nil {
		return "", err
	}
	if result.Error.Message != "" {
		return "", fmt.Errorf("failed to remove file: %s", result.Error.Message)
	}
	if result.Result != "ok" {
		return "", fmt.Errorf("failed to remove file: %s", result.Result)
	}
	return result.Result, nil
}
