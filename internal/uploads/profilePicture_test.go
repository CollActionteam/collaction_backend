package uploads_test

import (
	"context"
	"errors"
	"testing"

	"github.com/CollActionteam/collaction_backend/internal/uploads"
	"github.com/CollActionteam/collaction_backend/pkg/mocks/repository"
	"github.com/stretchr/testify/assert"
)

func TestProfile_GetUploadUrl(t *testing.T) {
	as := assert.New(t)
	imageUploadRepository := &repository.ProfilePicture{}

	t.Run("happy path", func(t *testing.T) {
		ext := "jpeg"
		userID := "1"
		service := uploads.NewProfileImageUploadService(imageUploadRepository)

		imageUploadRepository.On("GetUploadUrl", context.Background(), ext, userID).Return("https://sample-user-url-generated", nil).Once()
		url, err := service.GetUploadUrl(context.Background(), ext, userID)
		as.Equal("https://sample-user-url-generated", url)
		as.NoError(err)
		imageUploadRepository.AssertExpectations(t)
	})

	t.Run("input error: userID empty", func(t *testing.T) {
		ext := "jpeg"
		userID := ""
		service := uploads.NewProfileImageUploadService(imageUploadRepository)

		imageUploadRepository.On("GetUploadUrl", context.Background(), ext, userID).Return("", errors.New("user not found")).Once()
		url, err := service.GetUploadUrl(context.Background(), ext, userID)
		as.Equal("", url)
		as.Error(err)
		imageUploadRepository.AssertExpectations(t)
	})

	t.Run("input error: invalid image format", func(t *testing.T) {
		ext := "agsgfh"
		userID := "1"
		service := uploads.NewProfileImageUploadService(imageUploadRepository)

		imageUploadRepository.On("GetUploadUrl", context.Background(), ext, userID).Return("", errors.New("unknown file format")).Once()
		url, err := service.GetUploadUrl(context.Background(), ext, userID)
		as.Equal("", url)
		as.Error(err)
		imageUploadRepository.AssertExpectations(t)
	})

	t.Run("system error", func(t *testing.T) {
		ext := "jpeg"
		userID := "1"
		service := uploads.NewProfileImageUploadService(imageUploadRepository)

		imageUploadRepository.On("GetUploadUrl", context.Background(), ext, userID).Return("", errors.New("something failed")).Once()
		url, err := service.GetUploadUrl(context.Background(), ext, userID)
		as.Equal("", url)
		as.Error(err)
		imageUploadRepository.AssertExpectations(t)
	})
}
