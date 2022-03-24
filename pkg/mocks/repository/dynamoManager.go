package repository

import (
	models "github.com/CollActionteam/collaction_backend/internal/models"
	"github.com/CollActionteam/collaction_backend/utils"
	"github.com/stretchr/testify/mock"
)

type Dynamo struct {
	mock.Mock
}

func (d *Dynamo) GetDBItem(tableName string, pk string, sk string) (*models.CrowdactionData, error) {
	args := d.Called(tableName, pk, sk)
	return args.Get(0).(*models.CrowdactionData), args.Error(1)
}

func (d *Dynamo) Query(tableName string, filter string, startFrom *utils.PrimaryKey) ([]models.CrowdactionData, error) {
	args := d.Called(tableName, filter, startFrom)
	return args.Get(0).([]models.CrowdactionData), args.Error(1)
}
