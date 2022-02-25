package repository

import (
	"github.com/CollActionteam/collaction_backend/pkg/repository/aws"
	"github.com/stretchr/testify/mock"
)

type DynamoTable struct {
	mock.Mock
}

func (t *DynamoTable) DynamoGetItemKV(key, value string, receiver interface{}) error {
	args := t.Called(key, value, receiver)
	return args.Error(0)
}

func (t *DynamoTable) DynamoUpdateItemKV(data *aws.UpdateItemData) error {
	args := t.Called(data)
	return args.Error(0)
}

func (t *DynamoTable) DynamoInsertItemKV(data interface{}) error {
	args := t.Called(data)
	return args.Error(0)
}
