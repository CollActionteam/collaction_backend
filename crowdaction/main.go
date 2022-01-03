/* TODO Refactor!
 * - Move all functions relating to accessing the database to the models package.
 * - Move all constants to the utils package.
 */

package main

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/CollActionteam/collaction_backend/models"
	"github.com/CollActionteam/collaction_backend/utils"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

const (

	//do not send back the code/password, just an indication it's needed
	passwordRequired = "required"
)

var (
	tableName = os.Getenv("TABLE_NAME")
)

//get list of crowd actions
func getCrowdactions(status string) (events.APIGatewayProxyResponse, error) {
	var crowdactions []models.Crowdaction
	var err error
	switch status {
	case "joinable":
		crowdactions, _, err = models.ListJoinableCrowdactions(tableName, nil)
	case "active":
		crowdactions, _, err = models.ListActiveCrowdactions(tableName, nil)
	case "ended":
		crowdactions, _, err = models.ListCompletedCrowdactions(tableName, nil)
	default:
		return utils.GetMessageHttpResponse(http.StatusBadRequest, `unrecognized value for "status"`), nil
	}
	/* TODO Send password for handling in app for MVP
	for i; i < len(crowdactions); i++) {
		if crowdactions[i].PasswordJoin != "" {
			crowdactions[i].PasswordJoin = passwordRequired
		}
	}
	*/
	if err != nil {
		return utils.GetMessageHttpResponse(http.StatusInternalServerError, err.Error()), nil
	} else {
		body, err := json.Marshal(crowdactions)
		if err != nil {
			return utils.GetMessageHttpResponse(http.StatusInternalServerError, err.Error()), nil
		}
		return events.APIGatewayProxyResponse{
			Body:       string(body),
			StatusCode: http.StatusOK,
		}, nil
	}
}

//get details about a crowd action
func getCrowdaction(crowdactionID string) (events.APIGatewayProxyResponse, error) {
	crowdaction, err := models.GetCrowdaction(crowdactionID, tableName)
	if err != nil {
		return utils.GetMessageHttpResponse(http.StatusInternalServerError, err.Error()), nil
	}
	if crowdaction == nil {
		return utils.GetMessageHttpResponse(http.StatusNotFound, "crowdaction does not exist"), nil
	}
	/* TODO Send password for handling in app for MVP
	if crowdaction.PasswordJoin != "" {
		crowdaction.PasswordJoin = passwordRequired
	}
	*/
	body, err := json.Marshal(map[string]interface{}{"data": crowdaction})
	if err != nil {
		return utils.GetMessageHttpResponse(http.StatusInternalServerError, err.Error()), nil
	}
	return events.APIGatewayProxyResponse{
		Body:       string(body),
		StatusCode: http.StatusOK,
	}, nil

}

func handler(req events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
	crowdactionID := req.PathParameters["crowdactionID"]
	if crowdactionID == "" {
		status := req.QueryStringParameters["status"]
		switch status {
		case "":
			status = "joinable" //kind of default
		case "featured":
			status = "joinable" //for the time being
		}
		// TODO implement pagination
		return getCrowdactions(status)
	}
	return getCrowdaction(crowdactionID)
}

func main() {
	lambda.Start(handler)
}
