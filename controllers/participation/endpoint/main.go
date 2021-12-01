package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/CollActionteam/collaction_backend/auth"
	"github.com/CollActionteam/collaction_backend/models"
	"github.com/CollActionteam/collaction_backend/utils"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type JoinPayload struct {
	Password    string   `json:"password,omitempty"`
	Commitments []string `json:"commitments,omitempty"`
}

var (
	tableNameParticipation = os.Getenv("PARTICIPATION_TABLE")
	tableNameCrowdaction   = os.Getenv("CROWDACTION_TABLE")
	streamName             = os.Getenv("PARTICIPATION_STREAM")
)

func getParticipation(dbc *dynamodb.DynamoDB, userID string, crowdactionID string) (*models.ParticipationRecord, error) {
	out, err := dbc.GetItem(&dynamodb.GetItemInput{
		TableName: &tableNameParticipation,
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
	var r models.ParticipationRecord
	err = dynamodbattribute.UnmarshalMap(out.Item, r)
	return &r, err
}

func registerParticipation(userID string, name string, crowdaction *models.Crowdaction, payload JoinPayload) error {
	/* TODO Password not required when joining for MVP
	if crowdaction.PasswordJoin != "" && crowdaction.PasswordJoin != payload.Password {
		return fmt.Errorf("invalid password")
	}
	*/
	err := models.ValidateCommitments(payload.Commitments, crowdaction.CommitmentOptions)
	if err != nil {
		return err
	}
	sess := session.Must(session.NewSession())
	dbc := dynamodb.New(sess)
	part, err := getParticipation(dbc, userID, crowdaction.CrowdactionID)
	if part != nil {
		err = errors.New("already participating")
	}
	if err != nil {
		return err
	}
	av, err := dynamodbattribute.MarshalMap(models.ParticipationRecord{
		UserID:        userID,
		Name:          name,
		CrowdactionID: crowdaction.CrowdactionID,
		Commitments:   payload.Commitments,
		Date:          utils.GetDateStringNow(),
	})
	if err != nil {
		return err
	}
	_, err = dbc.PutItem(&dynamodb.PutItemInput{
		TableName: &tableNameParticipation,
		Item:      av,
	})
	if err == nil {
		err = recordEvent(sess, userID, crowdaction.CrowdactionID, payload.Commitments, +1)
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
		TableName: &tableNameParticipation,

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
		err = recordEvent(sess, userID, crowdactionID, part.Commitments, -1)
	}
	return err
}

func handler(req events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
	method := strings.ToLower(req.RequestContext.HTTP.Method)
	crowdactionID := req.PathParameters["crowdactionID"]
	var crowdaction models.Crowdaction
	usrInf, err := auth.ExtractUserInfo(req)
	if err != nil {
		crowdaction, err := models.GetCrowdaction(crowdactionID, tableNameCrowdaction)
		if err != nil {
			if !utils.IsFutureDateString(crowdaction.DateLimitJoin) {
				err = fmt.Errorf("cannot change participation for this crowdaction anymore")
			}
		}
	}
	if err == nil {
		if method == "post" {
			var payload JoinPayload
			err = json.Unmarshal([]byte(req.Body), &payload)
			if err == nil {
				err = registerParticipation(usrInf.UserID(), usrInf.Name(), &crowdaction, payload)
			}
		} else if method == "delete" {
			err = cancelParticipation(usrInf.UserID(), crowdactionID)
		} else if method == "get" {
			participation, err := getParticipation(utils.CreateDBClient(), usrInf.UserID(), crowdactionID)
			if err != nil {
				return utils.GetMessageHttpResponse(http.StatusInternalServerError, err.Error()), nil
			}
			var res events.APIGatewayProxyResponse
			if participation == nil {
				res = utils.GetMessageHttpResponse(http.StatusNotFound, "not participating")
			} else {
				// "Cannot go wrong"
				jsonPayload, _ := json.Marshal(participation)
				res = events.APIGatewayProxyResponse{
					Body:       string(jsonPayload),
					StatusCode: http.StatusOK,
				}
			}
			return res, nil
		} else {
			return utils.GetMessageHttpResponse(http.StatusNotImplemented, "not implemented"), nil
		}
	}
	if err != nil {
		return utils.GetMessageHttpResponse(http.StatusInternalServerError, err.Error()), nil
	} else {
		return utils.GetMessageHttpResponse(http.StatusOK, "updated"), nil
	}
}

func main() {
	lambda.Start(handler)
}
