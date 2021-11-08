package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
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
	//TO DO
	dbc := dynamodbClient()

	//https://dynobase.dev/dynamodb-golang-query-examples/#query

	//scan vs query

	_, err := dbc.Query(&dynamodb.QueryInput{
		TableName:              aws.String("CrowdactionTable"),
		KeyConditionExpression: aws.String("crowdactionID = :hashKey and start_datte > :rangeKey"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":hashKey":  &types.AttributeValueMemberS{Value: "123"},
			":rangeKey": &types.AttributeValueMemberN{Value: "20150101"},
		},
	})
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: http.StatusBadRequest}, nil
	}

	out, err := dbc.Scan(&dynamodb.ScanInput{
		TableName:        aws.String("CrowdactionTable"),
		FilterExpression: aws.String("attribute_not_exists(deletedAt) AND contains(firstName, :firstName)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":firstName": &types.AttributeValueMemberS{Value: "John"},
		},
	})
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: http.StatusBadRequest}, nil
	}

	fmt.Println(out.Items)

	//no error if we are here
	var body []byte
	return events.APIGatewayProxyResponse{
		Body:       string(body),
		StatusCode: http.StatusOK,
	}, nil
}

//get details about a crowd action
func getCrowdaction(crowdactionID string, req events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {

	dbc := dynamodbClient()

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

func dynamodbClient() *dynamodb.DynamoDB {
	sess := session.Must(session.NewSession())
	return dynamodb.New(sess)
}

func main() {
	lambda.Start(handler)
}
