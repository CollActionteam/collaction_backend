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
	tableName  = os.Getenv("TABLE_NAME")
	streamName = os.Getenv("PARTICIPATION_STREAM")
)

func getParticipation(dbClient *dynamodb.DynamoDB, userID string, crowdactionID string) (*models.ParticipationRecord, error) {
	pk := utils.PrefixPKparticipationUserID + userID
	sk := utils.PrefixSKparticipationCrowdactionID + crowdactionID
	out, err := utils.GetDBItem(dbClient, tableName, pk, sk)
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
	/* TODO replace with SQS
	if err == nil {
		err = recordEvent(sess, userID, crowdaction.CrowdactionID, payload.Commitments, +1)
	}
	*/
	return err
}

func cancelParticipation(userID string, crowdactionID string) error {
	dbClient := utils.CreateDBClient()
	part, err := getParticipation(dbClient, userID, crowdactionID)
	if err != nil {
		return err
	}
	if part == nil {
		return errors.New("not participating")
	}
	pk := utils.PrefixPKparticipationUserID + userID
	sk := utils.PrefixSKparticipationCrowdactionID + crowdactionID
	err = utils.DeleteDBItem(dbClient, tableName, pk, sk)
	/* TODO replace with SQS
	if err == nil {
		err = recordEvent(sess, userID, crowdactionID, part.Commitments, -1)
	}
	*/
	return err
}

func handler(req events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
	method := strings.ToLower(req.RequestContext.HTTP.Method)
	crowdactionID := req.PathParameters["crowdactionID"]
	var crowdaction models.Crowdaction
	usrInf, err := auth.ExtractUserInfo(req)
	if err != nil {
		crowdaction, err := models.GetCrowdaction(crowdactionID, tableName)
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
