package profile_test

import (
	"context"
	"testing"

	"github.com/CollActionteam/collaction_backend/internal/models"
	"github.com/CollActionteam/collaction_backend/internal/profile"
	"github.com/CollActionteam/collaction_backend/pkg/mocks/repository"
	"github.com/CollActionteam/collaction_backend/pkg/repository/aws"
	"github.com/stretchr/testify/assert"
)

func Test_UpdateProfile(t *testing.T) {
	ast := assert.New(t)
	profileRepository := &repository.Profile{
		UpdateData: aws.NewUpdateItem("searchKey", "searchValue", "updateDataKey", "updateDataValue"),
	}

	t.Run("testing UpdateProfile", func(t *testing.T) {
		userData := *models.NewUserInfo("user-id", "name", "+1993727843898")
		requestData := models.Profile{
			Country: "Canada",
			City:    "Calgary",
			Bio:     "My Bio",
		}

		profileRepository.On("UpdateUserProfile", context.Background(), userData, requestData).Return(nil).Once()

		service := profile.NewProfileCrudService(profileRepository)
		err := service.UpdateProfile(context.Background(), userData, requestData)
		ast.NoError(err)

		profileRepository.AssertExpectations(t)
	})
}

func Test_CreateProfile(t *testing.T) {
	ast := assert.New(t)
	dummy := models.Profile{
		DisplayName: "displayname",
		Country:     "Canada",
		City:        "Calgary",
		Bio:         "My Bio",
	}
	profileRepository := &repository.Profile{CreateData: dummy}

	t.Run("testing CreateProfile", func(t *testing.T) {
		userData := *models.NewUserInfo("user-id", "name", "+1993727843898")

		profileRepository.On("CreateUserProfile", context.Background(), userData, dummy).Return(nil).Once()

		service := profile.NewProfileCrudService(profileRepository)
		err := service.CreateProfile(context.Background(), userData, dummy)
		ast.NoError(err)

		profileRepository.AssertExpectations(t)
	})
}

func Test_GetProfile(t *testing.T) {
	ast := assert.New(t)
	testData := models.Profile{UserID: "user-id"}
	profileRepository := &repository.Profile{GetData: &testData}

	t.Run("testing GetProfile", func(t *testing.T) {
		userID := "user-id"
		service := profile.NewProfileCrudService(profileRepository)

		profileRepository.On("GetUserProfile", context.Background(), userID).Return(&testData, nil).Once()

		profileData, err := service.GetProfile(context.Background(), userID)
		ast.NoError(err)
		ast.Equal(&testData, profileData)

		profileRepository.AssertExpectations(t)
	})
}
