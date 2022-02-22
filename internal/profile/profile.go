package profile

import (
	"context"

	"github.com/CollActionteam/collaction_backend/internal/models"
)

type ProfileCrudRepository interface {
	GetUserProfile(ctx context.Context, userID string) (*models.Profile, error)
	UpdateUserProfile(ctx context.Context, user models.UserInfo, requestData models.Profile) error
	CreateUserProfile(ctx context.Context, user models.UserInfo, requestData models.Profile) error
}

type Service interface {
	GetProfile(ctx context.Context, userID string) (*models.Profile, error)
	UpdateProfile(ctx context.Context, user models.UserInfo, requestData models.Profile) error
	CreateProfile(ctx context.Context, user models.UserInfo, requestData models.Profile) error
}

type profile struct {
	profileRepository ProfileCrudRepository
}

func NewProfileCrudService(profileRepository ProfileCrudRepository) Service {
	return &profile{profileRepository: profileRepository}
}

func (p *profile) GetProfile(ctx context.Context, userID string) (*models.Profile, error) {
	return p.profileRepository.GetUserProfile(ctx, userID)
}

func (p *profile) UpdateProfile(ctx context.Context, user models.UserInfo, requestData models.Profile) error {
	return p.profileRepository.UpdateUserProfile(ctx, user, requestData)
}

func (p *profile) CreateProfile(ctx context.Context, user models.UserInfo, requestData models.Profile) error {
	return p.profileRepository.CreateUserProfile(ctx, user, requestData)
}
