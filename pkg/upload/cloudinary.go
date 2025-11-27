package upload

import (
	"context"
	"fmt"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/cloudinary/cloudinary-go/v2/logger"
	"github.com/rs/zerolog/log"
	"github.com/thanhphuocnguyen/go-eshop/config"
	applogger "github.com/thanhphuocnguyen/go-eshop/pkg/logger"
)

type CloudinaryUploader struct {
	uniqFileName bool
	cld          *cloudinary.Cloudinary
	cfg          config.Config
}

func NewCloudinaryUploader(cfg config.Config) CdnUploader {
	cld, err := cloudinary.NewFromURL(cfg.CloudinaryUrl)
	cld.Logger = &logger.Logger{
		Writer: applogger.NewLogger(nil),
	}

	if err != nil {
		panic(err)
	}

	return &CloudinaryUploader{
		cld:          cld,
		cfg:          cfg,
		uniqFileName: true,
	}
}

func (s *CloudinaryUploader) Upload(ctx context.Context, file interface{}) (publicID string, url string, err error) {
	useAssetFolderAsPublicIdPrefix := true
	result, err := s.cld.Upload.Upload(ctx, file, uploader.UploadParams{
		AssetFolder:                    s.cfg.CloudinaryFolder,
		UseAssetFolderAsPublicIDPrefix: &useAssetFolderAsPublicIdPrefix,
		UniqueFilename:                 &s.uniqFileName,
	})

	if err != nil {
		log.Err(err).Msg("failed to upload file")
		return "", "", err
	}

	publicID = result.PublicID
	url = result.SecureURL
	return
}

func (s *CloudinaryUploader) Remove(ctx context.Context, publicID string) (string, error) {
	// publicIDWithPath := fmt.Sprintf("%s/%s", s.cfg.CloudinaryFolder, publicID)
	result, err := s.cld.Upload.Destroy(ctx, uploader.DestroyParams{PublicID: publicID})
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
