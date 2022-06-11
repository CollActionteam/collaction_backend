package aws

import (
	"fmt"

	"github.com/CollActionteam/collaction_backend/internal/constants"
	m "github.com/CollActionteam/collaction_backend/internal/models"
	"github.com/CollActionteam/collaction_backend/utils"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

type Crowdaction interface {
	GetById(pk string, sk string) (*m.CrowdactionData, error)
	GetByStatus(status string, startFrom *utils.PrimaryKey) ([]m.CrowdactionData, error)
}

const (
	KeyDateStart      = "date_start"
	KeyDateEnd        = "date_end"
	KeyDateJoinBefore = "date_limit_join"
)

type crowdaction struct {
	dbClient *Dynamo
}

func NewCrowdaction(dynamo *Dynamo) Crowdaction {
	return &crowdaction{dbClient: dynamo}
}

/**
	GET Crowdaction by Id
**/
func (s *crowdaction) GetById(pk string, sk string) (*m.CrowdactionData, error) {
	item, err := s.dbClient.GetDBItem(constants.TableName, pk, sk)

	if item == nil || err != nil {
		return nil, err
	}

	var c m.CrowdactionData
	err = dynamodbattribute.UnmarshalMap(item, &c)

	return &c, err
}

/**
	GET Crowdaction by Status
**/
func (s *crowdaction) GetByStatus(status string, startFrom *utils.PrimaryKey) ([]m.CrowdactionData, error) {
	crowdactions := []m.CrowdactionData{}
	var filterCond expression.ConditionBuilder

	switch status {
	case "joinable":
		filterCond = expression.Name(KeyDateJoinBefore).GreaterThan(expression.Value(utils.GetDateStringNow()))
		fmt.Println("GetByStatus: joinable", filterCond)
	case "active":
		filterCond = expression.Name(KeyDateStart).LessThanEqual(expression.Value(utils.GetDateStringNow()))
		fmt.Println("GetByStatus: active", filterCond)
	case "ended":
		filterCond = expression.Name(KeyDateEnd).LessThanEqual(expression.Value(utils.GetDateStringNow()))
		fmt.Println("GetByStatus: ended", filterCond)
	default:
		fmt.Println("None of the edge cases matched")
	}

	items, err := s.dbClient.Query(constants.TableName, filterCond, startFrom)

	if items == nil || err != nil {
		return nil, err
	}

	for _, foo := range items {
		var crowdaction m.CrowdactionData
		err := dynamodbattribute.UnmarshalMap(foo, &crowdaction)

		if err == nil {
			crowdactions = append(crowdactions, crowdaction)
		}
	}

	if len(items) != len(crowdactions) {
		err = fmt.Errorf("error unmarshelling %d items", len(items)-len(crowdactions))
	}

	return crowdactions, err
}
