package repository

import (
	"context"

	models "github.com/CollActionteam/collaction_backend/internal/models"
	"github.com/CollActionteam/collaction_backend/utils"
	"github.com/stretchr/testify/mock"
)

type Dynamo struct {
	mock.Mock
}

func (d *Dynamo) GetCrowdactionById(ctx context.Context, crowdactionID string) error {
	args := d.Called(ctx, crowdactionID)
	return args.Error(0)
}

func (d *Dynamo) GetCrowdactionByStatus(ctx context.Context, status string, startFrom *utils.PrimaryKey) ([]models.CrowdactionData, error) {
	args := d.Called(ctx, status, startFrom)
	return args.Get(0).([]models.CrowdactionData), args.Error(1)
}
