package aws

import (
	"context"
	"strconv"

	"github.com/CollActionteam/collaction_backend/utils"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type crowdactionParticipations struct {
	table *DynamoTable
}

func NewCrowdactionParticipations(table *DynamoTable) *crowdactionParticipations {
	return &crowdactionParticipations{table: table}
}

func (e *crowdactionParticipations) ChangeCrowdactionParticipantCountBy(ctx context.Context, crowdactionID string, count int) error {
	dbClient := e.table.DbClient.dbClient
	_, err := dbClient.UpdateItem(&dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":c": {
				N: aws.String(strconv.Itoa(count)),
			},
		},
		TableName:        aws.String(e.table.Name),
		Key:              utils.GetPrimaryKey(utils.PKCrowdaction, crowdactionID),
		UpdateExpression: aws.String("set participant_count = participant_count + :c"),
	})
	return err
}
