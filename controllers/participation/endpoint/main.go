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
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type JoinPayload struct {
	Password    string   `json:"password,omitempty"`
	Commitments []string `json:"commitments,omitempty"`
}

var (
	tableName = os.Getenv("TABLE_NAME")
	queueUrl  = os.Getenv("PARTICIPATION_QUEUE")
)

func recordEvent(userID string, crowdactionID string, commitments []string, count int) error {
	qc := utils.CreateQueueClient()
	event := models.ParticipationEvent{
		UserID:        userID,
		CrowdactionID: crowdactionID,
		Commitments:   commitments,
		Count:         count,
	}
	return utils.SendQueueMessage(qc, queueUrl, event)
}

func getParticipation(dbClient *dynamodb.DynamoDB, userID string, crowdactionID string) (*models.ParticipationRecord, error) {
	pk := utils.PrefixPKparticipationUserID + userID
	sk := utils.PrefixSKparticipationCrowdactionID + crowdactionID
	item, err := utils.GetDBItem(dbClient, tableName, pk, sk)
	if item == nil || err != nil {
		return nil, err
	}
	var r models.ParticipationRecord
	err = dynamodbattribute.UnmarshalMap(item, &r)
	return &r, err
}

func registerParticipation(userID string, name string, crowdaction *models.Crowdaction, payload JoinPayload) error {
	if !utils.IsFutureDateString(crowdaction.DateLimitJoin) {
		return fmt.Errorf("cannot change participation for this crowdaction anymore")
	}
	/* TODO Password not required when joining for MVP
	if crowdaction.PasswordJoin != "" && crowdaction.PasswordJoin != payload.Password {
		return fmt.Errorf("invalid password")
	}
	*/
	err := models.ValidateCommitments(payload.Commitments, crowdaction.CommitmentOptions)
	if err != nil {
		return err
	}
	dbClient := utils.CreateDBClient()
	part, err := getParticipation(dbClient, userID, crowdaction.CrowdactionID)
	if part != nil {
		err = errors.New("already participating")
	}
	if err != nil {
		return err
	}
	pk := utils.PrefixPKparticipationUserID + userID
	sk := utils.PrefixSKparticipationCrowdactionID + crowdaction.CrowdactionID
	err = utils.PutDBItem(dbClient, tableName, pk, sk, models.ParticipationRecord{
		UserID:        userID,
		Name:          name,
		CrowdactionID: crowdaction.CrowdactionID,
		Commitments:   payload.Commitments,
		Date:          utils.GetDateStringNow(),
	})
	if err == nil {
		err = recordEvent(userID, crowdaction.CrowdactionID, payload.Commitments, +1)
	}
	return err
}

func cancelParticipation(userID string, crowdaction *models.Crowdaction) error {
	if !utils.IsFutureDateString(crowdaction.DateEnd) {
		return fmt.Errorf("cannot change participation for this crowdaction anymore")
	}
	dbClient := utils.CreateDBClient()
	part, err := getParticipation(dbClient, userID, crowdaction.CrowdactionID)
	if err != nil {
		return err
	}
	if part == nil {
		return errors.New("not participating")
	}
	pk := utils.PrefixPKparticipationUserID + userID
	sk := utils.PrefixSKparticipationCrowdactionID + crowdaction.CrowdactionID
	err = utils.DeleteDBItem(dbClient, tableName, pk, sk)
	if err == nil {
		err = recordEvent(userID, crowdaction.CrowdactionID, part.Commitments, -1)
	}
	return err
}

func handler(req events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
	method := strings.ToLower(req.RequestContext.HTTP.Method)
	crowdactionID := req.PathParameters["crowdactionID"]
	var crowdaction *models.Crowdaction
	usrInf, err := auth.ExtractUserInfo(req)
	if err != nil {
		return utils.GetMessageHttpResponse(http.StatusInternalServerError, err.Error()), nil
	}
	crowdaction, _ = models.GetCrowdaction(crowdactionID, tableName)
	if crowdaction == nil {
		return utils.GetMessageHttpResponse(http.StatusNotFound, "crowdaction not found"), nil
	}
	if method == "post" {
		var payload JoinPayload
		err = json.Unmarshal([]byte(req.Body), &payload)
		if err == nil {
			if len(payload.Commitments) == 0 {
				return utils.GetMessageHttpResponse(http.StatusBadRequest, "cannot participate without commitments"), nil
			}
			err = registerParticipation(usrInf.UserID(), usrInf.Name(), crowdaction, payload)
		}
	} else if method == "delete" {
		err = cancelParticipation(usrInf.UserID(), crowdaction)
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

	if err != nil {
		return utils.GetMessageHttpResponse(http.StatusInternalServerError, err.Error()), nil
	} else {
		return utils.GetMessageHttpResponse(http.StatusOK, "updated"), nil
	}
}

func main() {
	lambda.Start(handler)
}
