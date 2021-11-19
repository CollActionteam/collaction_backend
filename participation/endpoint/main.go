package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/CollActionteam/collaction_backend/auth"
	"github.com/CollActionteam/collaction_backend/participation"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/kinesis"
)

var (
	tableName  = os.Getenv("PARTICIPATION_TABLE")
	streamName = os.Getenv("PARTICIPATION_STREAM")
)

func getParticipation(dbc *dynamodb.DynamoDB, userID string, crowdactionID string) (*participation.ParticipationRecord, error) {
	out, err := dbc.GetItem(&dynamodb.GetItemInput{
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
	if out.Item == nil || err != nil {
		return nil, err
	}
	var r participation.ParticipationRecord
	err = dynamodbattribute.UnmarshalMap(out.Item, r)
	return &r, err
}

func recordEvent(sess *session.Session, userID string, crowdactionID string, commitmentID string, count int) error {
	kc := kinesis.New(sess)
	json, err := json.Marshal(participation.ParticipationEvent{
		UserID:        userID,
		CrowdactionID: crowdactionID,
		CommitmentID:  commitmentID,
		Count:         count,
	})
	if err != nil {
		return err
	}
	_, err = kc.PutRecord(&kinesis.PutRecordInput{
		StreamName:   &streamName,
		PartitionKey: &crowdactionID,
		Data:         json,
	})
	return err
}

func registerParticipation(userID string, name string, crowdactionID string, commitmentID string) error {
	// TODO check if commitmentID exists for crowdaction (Blocked by CAN-72)
	sess := session.Must(session.NewSession())
	dbc := dynamodb.New(sess)
	part, err := getParticipation(dbc, userID, crowdactionID)
	if part != nil {
		err = errors.New("already participating")
	}
	if err != nil {
		return err
	}
	av, err := dynamodbattribute.MarshalMap(participation.ParticipationRecord{
		UserID:        userID,
		Name:          name,
		CrowdactionID: crowdactionID,
		CommitmentID:  commitmentID,
		Timestamp:     time.Now().Unix(),
	})
	if err != nil {
		return err
	}
	_, err = dbc.PutItem(&dynamodb.PutItemInput{
		TableName: &tableName,
		Item:      av,
	})
	if err == nil {
		err = recordEvent(sess, userID, crowdactionID, commitmentID, +1)
	}
	return err
}

func cancelParticipation(userID string, crowdactionID string) error {
	sess := session.Must(session.NewSession())
	dbc := dynamodb.New(sess)
	part, err := getParticipation(dbc, userID, crowdactionID)
	if err != nil {
		return err
	}
	if part == nil {
		return errors.New("not participating")
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
		err = recordEvent(sess, userID, crowdactionID, part.CommitmentID, -1)
	}
	return err
}

func getResponseBody(msg string) string {
	// "Cannot go wrong"
	json, _ := json.Marshal(map[string]interface{}{"message": msg})
	return string(json)
}

func handler(req events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
	method := strings.ToLower(req.RequestContext.HTTP.Method)
	crowdactionID := req.PathParameters["crowdactionID"]
	// TODO check if crowdaction exists (Blocked by CAN-72)
	// TODO check if crowdaction is open for participation/leave (Blocked by CAN-72)
	usrInf, err := auth.ExtractUserInfo(req)
	if err == nil {
		if method == "post" {
			var payload participation.JoinPayload
			err = json.Unmarshal([]byte(req.Body), &payload)
			if err == nil {
				err = registerParticipation(usrInf.UserID(), usrInf.Name(), crowdactionID, payload.CommitmentID)
			}
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
			Body:       getResponseBody(err.Error()),
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
