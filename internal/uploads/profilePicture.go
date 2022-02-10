package uploads

type ProfileImageUploadRepository interface {
	GetUploadUrl(ext string, userID string) (string, error)
}

type Service interface {
	GetUploadUrl(ext string, userID string) (string, error)
}

type image struct {
	imageUploadRepository ProfileImageUploadRepository
}

func NewProfileImageUploadService(profileImageUploadRepo ProfileImageUploadRepository) Service {
	return &image{imageUploadRepository: profileImageUploadRepo}
}

func (i *image) GetUploadUrl(ext string, userID string) (string, error) {
	return i.imageUploadRepository.GetUploadUrl(ext, userID)
}
