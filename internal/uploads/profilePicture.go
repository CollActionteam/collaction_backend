package uploads

import "context"

type ProfileImageUploadRepository interface {
	GetUploadUrl(ctx context.Context, ext string, userID string) (string, error)
}

type Service interface {
	GetUploadUrl(ctx context.Context, ext string, userID string) (string, error)
}

type image struct {
	imageUploadRepository ProfileImageUploadRepository
}

func NewProfileImageUploadService(profileImageUploadRepo ProfileImageUploadRepository) Service {
	return &image{imageUploadRepository: profileImageUploadRepo}
}

func (i *image) GetUploadUrl(ctx context.Context, ext string, userID string) (string, error) {
	return i.imageUploadRepository.GetUploadUrl(ctx, ext, userID)
}
