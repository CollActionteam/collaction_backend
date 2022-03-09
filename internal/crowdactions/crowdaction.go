package crowdaction

import (
	"context"
	"fmt"

	"github.com/CollActionteam/collaction_backend/internal/constants"
	"github.com/CollActionteam/collaction_backend/internal/models"
	"github.com/CollActionteam/collaction_backend/utils"
)

const Separator = "### app version:"

type DynamoRepository interface {
	GetDBItem(tableName string, pk string, crowdactionId string) (*models.CrowdactionData, error)
}

type Service interface {
	GetCrowdaction(ctx context.Context, crowdactionId string) (*models.CrowdactionData, error)
}

type crowdaction struct {
	dynamodb DynamoRepository
}

func NewCrowdactionService(dynamodb DynamoRepository) Service {
	return &crowdaction{dynamodb: dynamodb}
}

/**
	GET Crowdaction by Id
**/
func (e *crowdaction) GetCrowdaction(ctx context.Context, crowdactionId string) (*models.CrowdactionData, error) {
	item, err := e.dynamodb.GetDBItem(constants.TableName, utils.PKCrowdaction, crowdactionId)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, fmt.Errorf("Crowdaction not found")
	}
	// var crowdaction models.CrowdactionData
	// err = dynamodbattribute.UnmarshalMap(foo, &crowdaction)
	return item, err
}
