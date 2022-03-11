package crowdaction

import (
	"context"
	"fmt"

	"github.com/CollActionteam/collaction_backend/internal/constants"
	"github.com/CollActionteam/collaction_backend/internal/models"
	"github.com/CollActionteam/collaction_backend/utils"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

const Separator = "### app version:"

type DynamoRepository interface {
	GetDBItem(tableName string, pk string, crowdactionId string) (map[string]*dynamodb.AttributeValue, error)
}

type Service interface {
	GetCrowdaction(ctx context.Context, crowdactionId string) (*models.CrowdactionData, error)
}

type crowdaction struct {
	dynamodb DynamoRepository
}

func NewCrowdactionService(dynamodb DynamoRepository) Service {
	// fmt.Println("NewCrowdactionService")
	return &crowdaction{dynamodb: dynamodb}
}

/**
	GET Crowdaction by Id
**/
func (e *crowdaction) GetCrowdaction(ctx context.Context, crowdactionID string) (*models.CrowdactionData, error) {
	// fmt.Println("GetCrowdaction calling internal:", crowdactionID)
	item, err := e.dynamodb.GetDBItem(constants.TableName, utils.PKCrowdaction, crowdactionID)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, fmt.Errorf("Crowdaction not found")
	}
	var crowdaction models.CrowdactionData
	err = dynamodbattribute.UnmarshalMap(item, &crowdaction)
	return &crowdaction, err
	// return item, err
}
