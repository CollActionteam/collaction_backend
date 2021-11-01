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

func handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	crowdactionID := req.PathParameters["crowdactionID"]

	sess := session.Must(session.NewSession())
	dbc := dynamodb.New(sess)

	out, err := dbc.GetItem(&dynamodb.GetItemInput{
		TableName: &tableName,
		Key: map[string]*dynamodb.AttributeValue{
			"crowdactionID": {
				S: aws.String(crowdactionID),
			},
		},
	})

	if err != nil {
		body, _ := json.Marshal(map[string]interface{}{"message": err.Error()})
		return events.APIGatewayProxyResponse{
			Body:       string(body),
			StatusCode: http.StatusInternalServerError,
		}, nil
	}

	if out.Item == nil {
		body, _ := json.Marshal(map[string]interface{}{"message": "crowdaction does not exist"})
		return events.APIGatewayProxyResponse{
			Body:       string(body),
			StatusCode: http.StatusNotFound,
		}, nil
	}

	var crowdaction Crowdaction
	dynamodbattribute.UnmarshalMap(out.Item, &crowdaction)

	body, _ := json.Marshal(map[string]interface{}{"data": crowdaction})
	return events.APIGatewayProxyResponse{
		Body:       string(body),
		StatusCode: http.StatusOK,
	}, nil
}

func main() {
	lambda.Start(handler)
}
