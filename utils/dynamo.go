package utils

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

const (
	//access pattern getCrowdaction
	//item has PK="act#<crowdactionID>" and SK="act#<crowdactionID>"
	PrefixPKcrowdactionID = "act#"
	PrefixSKcrowdactionID = "act#"

	//access pattern getActiveCrowdactions
	//item has PK="act_end#date_end" and SK="act#<crowdactionID>"
	PrefixPKcrowdaction_date_end = "act_end#"

	//access pattern getEligibleToJoinCrowdactions
	//item has PK="act_join#date_limit_join" and SK="date_start#act#<crowdactionID>"
	PrefixPKcrowdaction_date_limit_join = "act_join#"
)

func CreateDBClient() *dynamodb.DynamoDB {
	sess := session.Must(session.NewSession())
	return dynamodb.New(sess)
}

func GetDBItem(dbClient *dynamodb.DynamoDB, pk string, tableName string) (*dynamodb.GetItemOutput, error) {

	result, err := dbClient.GetItem(&dynamodb.GetItemInput{
		TableName: &tableName,
		Key: map[string]*dynamodb.AttributeValue{
			"pk": {
				S: aws.String(pk),
			},
			"sk": {
				S: aws.String(pk),
			},
		},
	})

	if err != nil {
		return nil, err
	}

	if result.Item == nil {
		msg := "could not find record '" + pk + "'"
		return nil, errors.New(msg)
	}

	return result, nil
}

func GetDBItems(dbClient *dynamodb.DynamoDB, pk string, sk string, tableName string) (*dynamodb.QueryOutput, error) {

	var keyCond expression.KeyConditionBuilder

	if sk == "" {
		keyCond = expression.Key("pk").Equal(expression.Value(pk))
	} else {
		keyCond = expression.KeyAnd(
			expression.Key("pk").Equal(expression.Value(pk)),
			expression.Key("sk").LessThan(expression.Value(sk)),
		)
	}

	expr, err := expression.NewBuilder().WithKeyCondition(keyCond).Build()
	if err != nil {
		return nil, err
	}

	input := &dynamodb.QueryInput{
		TableName:                 &tableName,
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
	}

	result, err := dbClient.Query(input)
	if err != nil {
		return nil, err
	}
	return result, nil
}
