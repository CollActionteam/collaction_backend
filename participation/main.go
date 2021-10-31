package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/CollActionteam/collaction_backend/auth"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/aws/aws-sdk-go/service/kinesis"
)

var (
	tableName  = os.Getenv("PARTICIPATION_TABLE")
	streamName = os.Getenv("PARTICIPATION_STREAM")
)

func doesParticipationExist(dbc *dynamodb.DynamoDB, userID string, crowdactionID string) (bool, error) {
	var err error
	keyCond := expression.KeyAnd(
		expression.Key("userID").Equal(expression.Value(userID)),
		expression.Key("crowdactionID").Equal(expression.Value(crowdactionID)),
	)
	expr, err := expression.NewBuilder().WithKeyCondition(keyCond).Build()
	if err != nil {
		return false, err
	}
	out, err := dbc.Query(&dynamodb.QueryInput{
		TableName:              &tableName,
		IndexName:              aws.String("participation"),
		KeyConditionExpression: expr.KeyCondition(),
	})
	exists := false
	if out != nil && out.Count != nil {
		exists = *out.Count != 1
	}
	return exists, err
}

func recordEvent(sess *session.Session, userID string, crowdactionID string, count int) error {
	kc := kinesis.New(sess)
	json, _ := json.Marshal(ParticipationEvent{
		UserID:        userID,
		CrowdactionID: crowdactionID,
		Count:         count,
	})
	_, err := kc.PutRecord(&kinesis.PutRecordInput{
		StreamName: &streamName,
		Data:       json,
	})
	return err
}

func registerParticipation(userID string, name string, crowdactionID string) error {
	var err error
	sess := session.Must(session.NewSession())
	dbc := dynamodb.New(sess)
	exists, err := doesParticipationExist(dbc, userID, crowdactionID)
	if exists {
		err = errors.New("already participating")
	}
	if err != nil {
		return err
	}
	av, err := dynamodbattribute.MarshalMap(map[string]interface{}{
		"userID":        userID,
		"name":          name,
		"crowdactionID": crowdactionID,
		"timestamp":     time.Now().Unix(),
	})
	if err != nil {
		return err
	}
	_, err = dbc.PutItem(&dynamodb.PutItemInput{
		TableName: &tableName,
		Item:      av,
	})
	if err == nil {
		err = recordEvent(sess, userID, crowdactionID, +1)
	}
	return err
}

func cancelParticipation(userID string, crowdactionID string) error {
	var err error
	sess := session.Must(session.NewSession())
	dbc := dynamodb.New(sess)
	exists, err := doesParticipationExist(dbc, userID, crowdactionID)
	if !exists {
		err = errors.New("not participating")
	}
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}
	_, err = dbc.DeleteItem(&dynamodb.DeleteItemInput{
		TableName: &tableName,
		Key: map[string]*dynamodb.AttributeValue{
			"userID": {
				S: aws.String(userID),
			},
			"crowdactionID": {
				S: aws.String(crowdactionID),
			},
		},
	})
	if err == nil {
		err = recordEvent(sess, userID, crowdactionID, -1)
	}
	return err
}

func getResponseBody(msg string) string {
	json, _ := json.Marshal(map[string]interface{}{"message": msg})
	return string(json)
}

func handler(req events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
	method := strings.ToLower(req.RequestContext.HTTP.Method)
	crowdactionID := req.PathParameters["crowdactionID"]
	var err error
	usrInf, err := auth.ExtractUserInfoV2(req)
	if err == nil {
		if method == "post" {
			err = registerParticipation(usrInf.UserID(), usrInf.Name(), crowdactionID)
		} else if method == "delete" {
			err = cancelParticipation(usrInf.UserID(), crowdactionID)
		} else {
			return events.APIGatewayProxyResponse{
				Body:       getResponseBody("Not implemented"),
				StatusCode: http.StatusNotImplemented,
			}, nil
		}
	}
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: http.StatusInternalServerError,
		}, nil
	} else {
		return events.APIGatewayProxyResponse{
			Body:       getResponseBody("updated"),
			StatusCode: http.StatusOK,
		}, nil
	}
}

func main() {
	lambda.Start(handler)
}
