package models

import (
	"fmt"
	"strconv"

	"github.com/CollActionteam/collaction_backend/utils"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

const (
	KeyDateStart      = "date_start"
	KeyDateEnd        = "date_end"
	KeyDateJoinBefore = "date_limit_join"
)

type CrowdactionParticipant struct {
	Name   string `json:"name,omitempty"`
	UserID string `json:"userID,omitempty"`
}

type CrowdactionImages struct {
	Card   string `json:"card,omitempty"`
	Banner string `json:"banner,omitempty"`
}

// TODO for later: Use different model between database and api (PasswordJoin vs. IsPasswordRequired)
type Crowdaction struct {
	CrowdactionID     string                   `json:"crowdactionID,omitempty"`
	Title             string                   `json:"title,omitempty"`
	Description       string                   `json:"description,omitempty"`
	Category          string                   `json:"category,omitempty"`
	Subcategory       string                   `json:"sub_category,omitempty"`
	Location          string                   `json:"location,omitempty"`
	DateStart         string                   `json:"date_start,omitempty"`      // Must match KeyDateStart
	DateEnd           string                   `json:"date_end,omitempty"`        // Must match KeyDateEnd
	DateLimitJoin     string                   `json:"date_limit_join,omitempty"` // Must match KeyDateLimitJoin
	PasswordJoin      string                   `json:"password_join,omitempty"`
	CommitmentOptions []CommitmentOption       `json:"commitment_options,omitempty"`
	ParticipantCount  int                      `json:"participant_count,omitempty"`
	TopParticipants   []CrowdactionParticipant `json:"top_participants,omitempty"`
	Images            CrowdactionImages        `json:"images,omitempty"`
}

func GetCrowdaction(crowdactionID string, tableName string) (*Crowdaction, error) {
	var crowdaction Crowdaction
	dbClient := utils.CreateDBClient()
	item, err := utils.GetDBItem(dbClient, tableName, utils.PKCrowdaction, crowdactionID)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, fmt.Errorf("crowdaction not found")
	}
	err = dynamodbattribute.UnmarshalMap(item, &crowdaction)
	if err != nil {
		return nil, err
	}
	return &crowdaction, nil
}

// TODO include parameter "skPrefix" to efficiently select category and subcategory
func listCrowdactions(tableName string, filterCond expression.ConditionBuilder, startFrom *utils.PrimaryKey) ([]Crowdaction, *utils.PrimaryKey, error) {
	crowdactions := []Crowdaction{}
	dbClient := utils.CreateDBClient()
	keyCond := expression.Key(utils.PartitionKey).Equal(expression.Value(utils.PKCrowdaction))
	expr, err := expression.NewBuilder().
		WithKeyCondition(keyCond).
		WithFilter(filterCond).
		Build()
	if err != nil {
		return crowdactions, nil, err
	}
	var exclusiveStartKey utils.PrimaryKey
	if startFrom != nil {
		exclusiveStartKey = *startFrom
	}
	out, err := dbClient.Query(&dynamodb.QueryInput{
		Limit:                     aws.Int64(utils.CrowdactionsPageLength),
		ExclusiveStartKey:         exclusiveStartKey,
		TableName:                 aws.String(tableName),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
		FilterExpression:          expr.Filter(),
	})
	if err != nil {
		return crowdactions, nil, err
	} else if out == nil {
		return crowdactions, nil, fmt.Errorf("query did not return result")
	}
	for _, item := range out.Items {
		var crowdaction Crowdaction
		itemErr := dynamodbattribute.UnmarshalMap(item, &crowdaction)
		if itemErr == nil {
			crowdactions = append(crowdactions, crowdaction)
		}
	}
	if len(out.Items) != len(crowdactions) {
		err = fmt.Errorf("error unmarshalling %d items", len(out.Items)-len(crowdactions))
	}
	startFrom = nil
	if out.LastEvaluatedKey != nil && len(out.LastEvaluatedKey) > 0 {
		lastPK := out.LastEvaluatedKey[utils.PartitionKey].S
		lastSK := out.LastEvaluatedKey[utils.SortKey].S
		exclusiveStartKey = utils.GetPrimaryKey(*lastPK, *lastSK)
		startFrom = &exclusiveStartKey
	}
	return crowdactions, startFrom, err
}

func ListActiveCrowdactions(tableName string, startFrom *utils.PrimaryKey) ([]Crowdaction, *utils.PrimaryKey, error) {
	filterCond := expression.Name(KeyDateStart).LessThanEqual(expression.Value(utils.GetDateStringNow()))
	return listCrowdactions(tableName, filterCond, startFrom)
}

func ListJoinableCrowdactions(tableName string, startFrom *utils.PrimaryKey) ([]Crowdaction, *utils.PrimaryKey, error) {
	filterCond := expression.Name(KeyDateJoinBefore).GreaterThan(expression.Value(utils.GetDateStringNow()))
	return listCrowdactions(tableName, filterCond, startFrom)
}

func ListCompletedCrowdactions(tableName string, startFrom *utils.PrimaryKey) ([]Crowdaction, *utils.PrimaryKey, error) {
	filterCond := expression.Name(KeyDateEnd).LessThanEqual(expression.Value(utils.GetDateStringNow()))
	return listCrowdactions(tableName, filterCond, startFrom)
}

func ChangeCrowdactionParticipantCountBy(crowdactionID string, tableName string, count int) error {
	dbClient := utils.CreateDBClient()
	_, err := dbClient.UpdateItem(&dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":c": {
				N: aws.String(strconv.Itoa(count)),
			},
		},
		TableName:        aws.String(tableName),
		Key:              utils.GetPrimaryKey(utils.PKCrowdaction, crowdactionID),
		UpdateExpression: aws.String("set participant_count = participant_count + :c"),
	})
	return err
}
