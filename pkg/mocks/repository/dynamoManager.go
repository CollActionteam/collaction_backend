package repository

import (
	"context"

	m "github.com/CollActionteam/collaction_backend/internal/models"
	"github.com/CollActionteam/collaction_backend/utils"
	"github.com/stretchr/testify/mock"
)

type Dynamo struct {
	mock.Mock
}

func (d *Dynamo) GetAll() ([]m.CrowdactionData, error) {
	args := d.Mock.Called()
	return args.Get(0).([]m.CrowdactionData), args.Error(1)
}

func (d *Dynamo) GetById(pk string, sk string) (*m.CrowdactionData, error) {
	args := d.Mock.Called(pk, sk)
	return args.Get(0).(*m.CrowdactionData), args.Error(1)
}

func (d *Dynamo) GetByStatus(filterCond string, startFrom *utils.PrimaryKey) ([]m.CrowdactionData, error) {
	args := d.Mock.Called(filterCond, startFrom)
	return args.Get(0).([]m.CrowdactionData), args.Error(1)
}

func (d *Dynamo) Register(ctx context.Context, payload m.CrowdactionData) (*m.CrowdactionData, error) {
	d.Mock.Called(payload)
	return &payload, nil
}
