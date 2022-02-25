package repository

import (
	"context"

	"github.com/CollActionteam/collaction_backend/internal/models"
	"github.com/CollActionteam/collaction_backend/pkg/repository/aws"
	"github.com/stretchr/testify/mock"
)

type Profile struct {
	mock.Mock
	UpdateData *aws.UpdateItemData
	CreateData interface{}
	GetData    *models.Profile
}

func (p *Profile) GetUserProfile(ctx context.Context, userID string) (*models.Profile, error) {
	dynamoClient, receiver := DynamoTable{}, models.Profile{}
	dynamoClient.On("DynamoGetItemKV", "key", "value", receiver).Return(nil).Once()

	args := p.Called(ctx, userID)
	return p.GetData, args.Error(1)
}

func (p *Profile) UpdateUserProfile(ctx context.Context, user models.UserInfo, requestData models.Profile) error {
	dynamoClient := DynamoTable{}
	dynamoClient.On("DynamoUpdateItemKV", p.UpdateData).Return(nil).Once()

	args := p.Called(ctx, user, requestData)
	return args.Error(0)
}

func (p *Profile) CreateUserProfile(ctx context.Context, user models.UserInfo, requestData models.Profile) error {
	dynamoClient := DynamoTable{}
	dynamoClient.On("DynamoInsertItemKV", p.CreateData).Return(nil).Once()

	args := p.Called(ctx, user, requestData)
	return args.Error(0)
}
