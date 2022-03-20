package crowdaction

import (
	"context"
	"fmt"

	"github.com/CollActionteam/collaction_backend/internal/constants"
	"github.com/CollActionteam/collaction_backend/internal/models"
	"github.com/CollActionteam/collaction_backend/utils"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

type DynamoRepository interface {
	GetDBItem(tableName string, pk string, crowdactionId string) (map[string]*dynamodb.AttributeValue, error)
	Query(tableName string, filterCond expression.ConditionBuilder, startFrom *utils.PrimaryKey) ([]models.CrowdactionData, error)
}
type Service interface {
	GetCrowdactionById(ctx context.Context, crowdactionId string) (*models.CrowdactionData, error)
	GetCrowdactionsByStatus(ctx context.Context, status string, startFrom *utils.PrimaryKey) ([]models.CrowdactionData, error)
}
type crowdaction struct {
	dynamodb DynamoRepository
}

const (
	KeyDateStart      = "date_start"
	KeyDateEnd        = "date_end"
	KeyDateJoinBefore = "date_limit_join"
)

func NewCrowdactionService(dynamodb DynamoRepository) Service {
	return &crowdaction{dynamodb: dynamodb}
}

/**
	GET Crowdaction by Id
**/
func (e *crowdaction) GetCrowdactionById(ctx context.Context, crowdactionID string) (*models.CrowdactionData, error) {
	fmt.Println("GetCrowdaction calling internal:", crowdactionID)
	item, err := e.dynamodb.GetDBItem(constants.TableName, utils.PKCrowdaction, crowdactionID)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, fmt.Errorf("crowdaction not found")
	}
	var crowdaction models.CrowdactionData
	err = dynamodbattribute.UnmarshalMap(item, &crowdaction)
	return &crowdaction, err
}

/**
	GET Crowdaction by Status
**/
func (e *crowdaction) GetCrowdactionsByStatus(ctx context.Context, status string, startFrom *utils.PrimaryKey) ([]models.CrowdactionData, error) {
	crowdactions := []models.CrowdactionData{} // empty crowdaction array

	switch status {
	case "joinable":
		filterCond := expression.Name(KeyDateJoinBefore).GreaterThan(expression.Value(utils.GetDateStringNow()))
		items, err := e.dynamodb.Query(constants.TableName, filterCond, startFrom)
		return items, err
	case "active":
		filterCond := expression.Name(KeyDateStart).LessThanEqual(expression.Value(utils.GetDateStringNow()))
		items, err := e.dynamodb.Query(constants.TableName, filterCond, startFrom)
		return items, err
	case "ended":
		filterCond := expression.Name(KeyDateEnd).LessThanEqual(expression.Value(utils.GetDateStringNow()))
		items, err := e.dynamodb.Query(constants.TableName, filterCond, startFrom)
		return items, err
	default:
		return crowdactions, nil
	}
}
