/* TODO Refactor!
 * - Move all functions relating to accessing the database to the models package.
 * - Move all constants to the utils package.
 */

package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/CollActionteam/collaction_backend/models"
	"github.com/CollActionteam/collaction_backend/utils"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

const (

	//do not send back the code/password, just an indication it's needed
	passwordRequired = "required"

	//for the time being let's consider 30 days how far in future we retrieve crowdactions
	//could use a parameter store here
	//example: current date = 2021.12.01 => retrieve up to and including crowdaction end date=2021.12.30
	dateLimit = 30
)

type pkError struct {
	Pk  string `json:"pk"`
	Err string `json:"error"`
}

var (
	tableName = os.Getenv("TABLE_NAME")
	dbClient  *dynamodb.DynamoDB
)

//get list of crowd actions
func getListCrowdaction(req events.APIGatewayV2HTTPRequest, status string) (events.APIGatewayProxyResponse, error) {

	var pk, sk string
	var items []*dynamodb.QueryOutput
	var crowdaction models.Crowdaction
	var errors []string

	mapErr := make(map[string]error)

	//get current date
	dateCurrent := time.Now()

	for i := 0; i < dateLimit; i++ {

		date := dateCurrent.AddDate(0, 0, i).Format(utils.DateFormat)
		switch status {
		case "active":
			pk = utils.PrefixPKcrowdaction_date_end + date
			sk = ""
		case "joinable":
			pk = utils.PrefixPKcrowdaction_date_limit_join + date
			//add one day so the result includes the crowdactions starting on the current day as well
			sk = dateCurrent.AddDate(0, 0, 1).Format(utils.DateFormat)
		}
		//get items for a partition key
		result, err := utils.GetDBItems(dbClient, pk, sk, tableName)
		//error handling https://play.golang.org/p/D5xTeZP9VnU
		if err != nil {
			mapErr[pk] = err
		}
		if result != nil {
			//add partition key result to the total result for the date range
			items = append(items, result)
		}
	}

	if len(mapErr) > 0 {
		e := pkError{}
		for k, v := range mapErr {
			e.Pk = k
			e.Err = v.Error()
			jsonErr, err := json.Marshal(e)
			if err != nil {
				errors = append(errors, err.Error())
			} else {
				errors = append(errors, string(jsonErr))
			}
		}

		body := "[" + strings.Join(errors, ",") + "]"
		return events.APIGatewayProxyResponse{Body: body, StatusCode: http.StatusInternalServerError}, nil

	}

	if items == nil {
		body, _ := json.Marshal(map[string]string{"message": "no active crowdactions found"})
		return events.APIGatewayProxyResponse{
			Body:       string(body),
			StatusCode: http.StatusNotFound,
		}, nil
	}

	//no error yet if we are here
	var listCrowdaction []string
	for _, v := range items {
		for _, item := range v.Items {
			err := dynamodbattribute.UnmarshalMap(item, &crowdaction)
			if err != nil {
				return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: http.StatusBadRequest}, nil
			}
			/* TODO Send password for handling in app for MVP
			if crowdaction.PasswordJoin != "" {
				crowdaction.PasswordJoin = passwordRequired
			}
			*/
			action, err := json.Marshal(map[string]interface{}{"data": crowdaction})
			if err != nil {
				return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: http.StatusBadRequest}, nil
			}
			listCrowdaction = append(listCrowdaction, string(action))
		}
	}
	body := "[" + strings.Join(listCrowdaction, ",") + "]"
	return events.APIGatewayProxyResponse{
		Body:       body,
		StatusCode: http.StatusOK,
	}, nil
}

//get details about a crowd action
func getCrowdaction(crowdactionID string, req events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {

	var crowdaction models.Crowdaction

	k := utils.PrefixPKcrowdactionID + crowdactionID
	item, err := utils.GetDBItem(dbClient, tableName, k, k)

	if err != nil {
		body, err := json.Marshal(map[string]interface{}{"message": err.Error()})
		if err != nil {
			return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: http.StatusBadRequest}, nil
		}
		return events.APIGatewayProxyResponse{
			Body:       string(body),
			StatusCode: http.StatusInternalServerError,
		}, nil
	}

	if item == nil {
		body, err := json.Marshal(map[string]string{"message": "crowdaction does not exist"})
		if err != nil {
			return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: http.StatusBadRequest}, nil
		}
		return events.APIGatewayProxyResponse{
			Body:       string(body),
			StatusCode: http.StatusNotFound,
		}, nil
	}

	err = dynamodbattribute.UnmarshalMap(item, &crowdaction)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: http.StatusBadRequest}, nil
	}
	/* TODO Send password for handling in app for MVP
	if crowdaction.PasswordJoin != "" {
		crowdaction.PasswordJoin = passwordRequired
	}
	*/
	body, err := json.Marshal(map[string]interface{}{"data": crowdaction})
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: http.StatusBadRequest}, nil
	}
	return events.APIGatewayProxyResponse{
		Body:       string(body),
		StatusCode: http.StatusOK,
	}, nil

}

func handler(req events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {

	var (
		resp events.APIGatewayProxyResponse
		err  error
	)

	dbClient = utils.CreateDBClient()

	crowdactionID := req.PathParameters["crowdactionID"]

	if crowdactionID == "" {
		status := req.QueryStringParameters["status"]

		switch status {
		case "":
			status = "joinable" //kind of default
		case "featured":
			status = "joinable" //for the time being
		case "joinable", "active", "ended":
		default:
			err := errors.New("unrecognizable status value")
			return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: http.StatusBadRequest}, nil
		}

		resp, err = getListCrowdaction(req, status)
		return resp, err
	}

	resp, err = getCrowdaction(crowdactionID, req)
	return resp, err
}

func main() {
	lambda.Start(handler)
}
