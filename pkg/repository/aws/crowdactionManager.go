package aws

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/CollActionteam/collaction_backend/internal/constants"
	m "github.com/CollActionteam/collaction_backend/internal/models"
	"github.com/CollActionteam/collaction_backend/utils"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

type Crowdaction interface {
	GetAll() ([]m.CrowdactionData, error)
	GetById(pk string, sk string) (*m.CrowdactionData, error)
	GetByStatus(status string, startFrom *utils.PrimaryKey) ([]m.CrowdactionData, error)
	Register(ctx context.Context, payload m.CrowdactionData) (*m.CrowdactionData, error)
	// Register(ctx context.Context, payload m.CrowdactionData) error
}

const (
	KeyDateStart      = "date_start"
	KeyDateEnd        = "date_end"
	KeyDateJoinBefore = "date_limit_join"
)

const charset = "abcdefghijklmnopqrstuvwxyz0123456789"

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

type crowdaction struct {
	dbClient *Dynamo
}

func NewCrowdaction(dynamo *Dynamo) Crowdaction {
	return &crowdaction{dbClient: dynamo}
}

func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func RandomIDPrefix(length int) string {
	return StringWithCharset(length, charset)
}

func (s *crowdaction) GetById(pk string, sk string) (*m.CrowdactionData, error) {
	item, err := s.dbClient.GetDBItem(constants.TableName, pk, sk)

	if item == nil || err != nil {
		return nil, err
	}

	var c m.CrowdactionData
	err = dynamodbattribute.UnmarshalMap(item, &c)

	return &c, err
}

func (s *crowdaction) GetAll() ([]m.CrowdactionData, error) {
	crowdactions := []m.CrowdactionData{} // crowdactions array
	var filterCond = expression.Name(utils.PartitionKey).Equal(expression.Value(utils.PKCrowdaction))

	item, err := s.dbClient.Scan(constants.TableName, filterCond)
	if item == nil || err != nil {
		return nil, err
	}

	for _, itemIterator := range item {
		var crowdaction m.CrowdactionData
		err := dynamodbattribute.UnmarshalMap(itemIterator, &crowdaction)

		if err == nil {
			crowdactions = append(crowdactions, crowdaction)
		}
	}

	if len(item) != len(crowdactions) {
		err = fmt.Errorf("error unmarshallaing %d items", len(item)-len(crowdactions))
	}

	return crowdactions, err
}

func (s *crowdaction) GetByStatus(status string, startFrom *utils.PrimaryKey) ([]m.CrowdactionData, error) {
	crowdactions := []m.CrowdactionData{}
	var filterCond expression.ConditionBuilder

	switch status {
	case "joinable":
		filterCond = expression.Name(KeyDateJoinBefore).GreaterThan(expression.Value(utils.GetDateStringNow()))
	case "active":
		filterCond = expression.Name(KeyDateStart).LessThanEqual(expression.Value(utils.GetDateStringNow()))
	case "ended":
		filterCond = expression.Name(KeyDateEnd).LessThanEqual(expression.Value(utils.GetDateStringNow()))
	default:
	}

	items, err := s.dbClient.Query(constants.TableName, filterCond, startFrom)

	if items == nil || err != nil {
		return nil, err
	}

	for _, itemIterator := range items {
		var crowdaction m.CrowdactionData
		err := dynamodbattribute.UnmarshalMap(itemIterator, &crowdaction)

		if err == nil {
			crowdactions = append(crowdactions, crowdaction)
		}
	}

	if len(items) != len(crowdactions) {
		err = fmt.Errorf("error unmarshelling %d items", len(items)-len(crowdactions))
	}

	return crowdactions, err
}

func (s *crowdaction) Register(ctx context.Context, payload m.CrowdactionData) (*m.CrowdactionData, error) {
	var response m.CrowdactionData
	generatedID := RandomIDPrefix(8)
	pk := utils.PKCrowdaction
	sk := payload.Category + "#" + payload.Subcategory + "#" + generatedID
	payload.CrowdactionID = sk
	response = payload
	// should modify the payload here to include the crowdcationID
	// fmt.Println("payload", payload)
	fmt.Println("9. pkg/respository/aws/crowdactionManager.go", payload)

	err := s.dbClient.PutDBItem(constants.TableName, pk, sk, payload)

	if err != nil {
		return &response, nil
	}
	return &response, err
}
