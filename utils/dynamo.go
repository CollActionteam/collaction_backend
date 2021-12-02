package utils

import (
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
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

	//access pattern getParticipation
	//item has PK="prt#<userID>" and SK="prt#<crowdactionID>"
	//(we want strong consistency when listing the users participation)
	PrefixPKparticipationUserID        = "prt#" + "usr#" // TODO refactor
	PrefixSKparticipationCrowdactionID = "prt" + PrefixPKcrowdactionID
)

func CreateDBClient() *dynamodb.DynamoDB {
	sess := session.Must(session.NewSession())
	return dynamodb.New(sess)
}

func GetDBItem(dbClient *dynamodb.DynamoDB, tableName string, pk string, sk string) (*dynamodb.GetItemOutput, error) {
	result, err := dbClient.GetItem(&dynamodb.GetItemInput{
		TableName: &tableName,
		Key: map[string]*dynamodb.AttributeValue{
			"pk": {
				S: aws.String(pk),
			},
			"sk": {
				S: aws.String(sk),
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

func PutDBItem(dbClient *dynamodb.DynamoDB, tableName string, pk string, sk string, record interface{}) error {
	av, err := dynamodbattribute.MarshalMap(record)
	if _, hasKey := av["pk"]; hasKey {
		return fmt.Errorf("record must not have a field with the label \"pk\"")
	}
	if _, hasKey := av["sk"]; hasKey {
		return fmt.Errorf("record must not have a field with the label \"sk\"")
	}
	*av["pk"] = dynamodb.AttributeValue{S: aws.String(pk)}
	*av["sk"] = dynamodb.AttributeValue{S: aws.String(sk)}
	if err != nil {
		_, err = dbClient.PutItem(&dynamodb.PutItemInput{
			TableName: &tableName,
			Item:      av,
		})
	}
	return err
}

func DeleteDBItem(dbClient *dynamodb.DynamoDB, tableName string, pk string, sk string) error {
	_, err := dbClient.DeleteItem(&dynamodb.DeleteItemInput{
		TableName: &tableName,
		Key: map[string]*dynamodb.AttributeValue{
			"pk": {
				S: aws.String(pk),
			},
			"sk": {
				S: aws.String(sk),
			},
		},
	})
	return err
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
