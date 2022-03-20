package repository

import (
	"context"

	"github.com/CollActionteam/collaction_backend/internal/models"
	"github.com/stretchr/testify/mock"
)

type Dynamo struct {
	mock.Mock
}

func (d *Dynamo) GetCrowdactionById(ctx context.Context, data models.CrowdactionData) error {
	args := d.Called(ctx, data)
	return args.Error(0)
}
