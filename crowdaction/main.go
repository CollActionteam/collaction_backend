package main

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

var (
	tableName = os.Getenv("CROWDACTION_TABLE")
)

////get list of crowd actions
func getListCrowdaction(req events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {

	//no error
	var body []byte
	return events.APIGatewayProxyResponse{
		Body:       string(body),
		StatusCode: http.StatusOK,
	}, nil
}

//get details about a crowd action
func getCrowdaction(crowdactionID string, req events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {

	dbc := getDBConnection()

	out, err := dbc.GetItem(&dynamodb.GetItemInput{
		TableName: &tableName,
		Key: map[string]*dynamodb.AttributeValue{
			"crowdactionID": {
				S: aws.String(crowdactionID),
			},
		},
	})

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

	if out.Item == nil {
		body, err := json.Marshal(map[string]interface{}{"message": "crowdaction does not exist"})
		if err != nil {
			return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: http.StatusBadRequest}, nil
		}
		return events.APIGatewayProxyResponse{
			Body:       string(body),
			StatusCode: http.StatusNotFound,
		}, nil
	}

	var crowdaction Crowdaction
	dynamodbattribute.UnmarshalMap(out.Item, &crowdaction)

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

	crowdactionID := req.PathParameters["crowdactionID"]

	if crowdactionID == "" {
		resp, err = getListCrowdaction(req)
		return resp, err
	}

	resp, err = getCrowdaction(crowdactionID, req)
	return resp, err
}

func getDBConnection() *dynamodb.DynamoDB {
	sess := session.Must(session.NewSession())
	return dynamodb.New(sess)
}

func main() {
	lambda.Start(handler)
}
