package repository

import (
	models "github.com/CollActionteam/collaction_backend/internal/models"
	"github.com/CollActionteam/collaction_backend/utils"
	"github.com/stretchr/testify/mock"
)

type Dynamo struct {
	mock.Mock
}

func (d *Dynamo) GetById(pk string, sk string) (*models.CrowdactionData, error) {
	args := d.Called(pk, sk)
	return args.Get(0).(*models.CrowdactionData), args.Error(1)
}

func (d *Dynamo) GetByStatus(filterCond string, startFrom *utils.PrimaryKey) ([]models.CrowdactionData, error) {
	args := d.Called(filterCond, startFrom)
	return args.Get(0).([]models.CrowdactionData), args.Error(1)
}
