package upload

import (
	"context"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/thanhphuocnguyen/go-eshop/config"
)

type CloudinaryUploadService struct {
	cld *cloudinary.Cloudinary
	cfg config.Config
}

func NewCloudinaryUploadService(cfg config.Config) UploadService {
	cld, err := cloudinary.NewFromURL(cfg.CloudinaryUrl)
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

func (s *CloudinaryUploadService) RemoveFile(ctx context.Context, filename string) error {
	_, err := s.cld.Upload.Destroy(ctx, uploader.DestroyParams{PublicID: filename})
	if err != nil {
		return err
	}
	return nil
}
