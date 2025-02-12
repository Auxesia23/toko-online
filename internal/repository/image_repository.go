package repository

import (
	"context"
	"mime/multipart"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type ImageRepository interface {
	Upload(ctx context.Context, file multipart.File, name string) (string, error)
}

type ImageRepo struct {
	CLD *cloudinary.Cloudinary
}

func NewImageRepository(cld *cloudinary.Cloudinary) ImageRepository {
	return &ImageRepo{
		CLD: cld,
	}
}

func (repo *ImageRepo) Upload(ctx context.Context, file multipart.File, name string) (string, error) {
	uploadResult, err := repo.CLD.Upload.Upload(ctx, file, uploader.UploadParams{
		PublicID: name,        
	})
	if err != nil {
		return "", err
	}

	return uploadResult.SecureURL, nil 
}
