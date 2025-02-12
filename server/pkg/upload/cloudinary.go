package upload

import (
	"context"
	"fmt"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/cloudinary/cloudinary-go/v2/logger"
	"github.com/thanhphuocnguyen/go-eshop/config"
	applogger "github.com/thanhphuocnguyen/go-eshop/pkg/logger"
)

type CloudinaryUploadService struct {
	uniqFileName bool
	cld          *cloudinary.Cloudinary
	cfg          config.Config
}

func NewCloudinaryUploadService(cfg config.Config) UploadService {
	cld, err := cloudinary.NewFromURL(cfg.CloudinaryUrl)
	cld.Logger = &logger.Logger{
		Writer: applogger.NewLogger(nil),
	}

	if err != nil {
		panic(err)
	}

	return &CloudinaryUploadService{
		cld:          cld,
		cfg:          cfg,
		uniqFileName: true,
	}
}

func (s *CloudinaryUploadService) UploadFile(ctx context.Context, file interface{}) (publicID string, url string, err error) {
	result, err := s.cld.Upload.Upload(ctx, file, uploader.UploadParams{
		AssetFolder:    s.cfg.CloudinaryFolder,
		UniqueFilename: &s.uniqFileName,
	})

	if err != nil {
		return "", "", err
	}

	publicID = result.PublicID
	url = result.SecureURL
	return
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
