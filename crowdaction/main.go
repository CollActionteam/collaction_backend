package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

const (

	//date format
	dateFormat = "20060102"

	//for the time being let's consider 30 days how far in future we retrieve crowdactions
	//could use a parameter store here
	//example: current date = 2021.12.01 => retrieve up to and including crowdaction end date=2021.12.30
	dateLimit = 30

	//access pattern getCrowdaction
	//item has PK="act#<crowdactionID>" and SK="act#<crowdactionID>"
	prefixPKcrowdactionID = "act#"
	prefixSKcrowdactionID = "act#"

	//access pattern getActiveCrowdactions
	//item has PK="act_end#date_end" and SK="act#<crowdactionID>"
	prefixPKcrowdaction_date_end = "act_end#"

	//access pattern getEligibleToJoinCrowdactions
	//item has PK="act_join#date_limit_join" and SK="date_start#act#<crowdactionID>"
	prefixPKcrowdaction_date_limit_join = "act_join#"
)

type pkError struct {
	Pk  string `json:"pk"`
	Err string `json:"error"`
}

var (
	tableName = os.Getenv("CROWDACTION_TABLE")
	dbClient  *dynamodb.DynamoDB
)

////get list of crowd actions
func getListCrowdaction(req events.APIGatewayV2HTTPRequest, status string) (events.APIGatewayProxyResponse, error) {

	var pk, sk string
	var items []*dynamodb.QueryOutput
	var crowdaction Crowdaction
	var errors []string

	mapErr := make(map[string]error)

	//get current date
	dateCurrent := time.Now()

	for i := 0; i < dateLimit; i++ {

		date := dateCurrent.AddDate(0, 0, i).Format(dateFormat)
		switch status {
		case "active":
			pk = prefixPKcrowdaction_date_end + date
			sk = ""
		case "joinable":
			pk = prefixPKcrowdaction_date_limit_join + date
			sk = dateCurrent.Format(dateFormat)
		}
		//get items for a partition key
		result, err := getItems(pk, sk)
		//error handling https://play.golang.org/p/D5xTeZP9VnU
		if err != nil {
			mapErr[pk] = err
		}
		if result != nil {
			//add partition key result to the total result for the date range
			items = append(items, result)
		}
	}

	if mapErr != nil {
		e := pkError{}
		for k, v := range mapErr {
			e.Pk = k
			e.Err = v.Error()
			jsonErr, _ := json.Marshal(e)
			errors = append(errors, string(jsonErr))
		}

		body := "[" + strings.Join(errors, ",") + "]"
		return events.APIGatewayProxyResponse{Body: body, StatusCode: http.StatusBadRequest}, nil
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
			c, err := json.Marshal(map[string]interface{}{"data": crowdaction})
			if err != nil {
				return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: http.StatusBadRequest}, nil
			}
			listCrowdaction = append(listCrowdaction, string(c))
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

	var crowdaction Crowdaction

	val := prefixPKcrowdactionID + crowdactionID
	out, err := getItem(val)

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
		body, err := json.Marshal(map[string]string{"message": "crowdaction does not exist"})
		if err != nil {
			return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: http.StatusBadRequest}, nil
		}
		return events.APIGatewayProxyResponse{
			Body:       string(body),
			StatusCode: http.StatusNotFound,
		}, nil
	}

	err = dynamodbattribute.UnmarshalMap(out.Item, &crowdaction)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: http.StatusBadRequest}, nil
	}

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

	dbClient = dynamodbClient()

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
			return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 400}, nil
		}

		resp, err = getListCrowdaction(req, status)
		return resp, err
	}

	resp, err = getCrowdaction(crowdactionID, req)
	return resp, err
}

func dynamodbClient() *dynamodb.DynamoDB {
	sess := session.Must(session.NewSession())
	return dynamodb.New(sess)
}

func getItem(val string) (*dynamodb.GetItemOutput, error) {

	result, err := dbClient.GetItem(&dynamodb.GetItemInput{
		TableName: &tableName,
		Key: map[string]*dynamodb.AttributeValue{
			"pk": {
				S: aws.String(val),
			},
			"sk": {
				S: aws.String(val),
			},
		},
	})

	if err != nil {
		return nil, err
	}

	if result.Item == nil {
		msg := "could not find crowdaction '" + val + "'"
		return nil, errors.New(msg)
	}

	return result, nil
}

func getItems(pk, sk string) (*dynamodb.QueryOutput, error) {

	if sk == "" {
		keyCond := expression.KeyAnd(
			expression.Key("pk").Equal(expression.Value(pk))
		)
	} else {
		keyCond := expression.KeyAnd(
			expression.Key("pk").Equal(expression.Value(pk)),
			expression.Key("sk").LessThanEqual(expression.Value(sk)),
		)
	}

	expr, err := expression.NewBuilder().WithKeyCondition(keyCond).Build()
	if err != nil {
		return nil, err
	}

	/*
			input := &dynamodb.QueryInput{
				ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
					":action_end": {S: aws.String(pk)},
				},
				KeyConditionExpression: aws.String("pk = :action_end"), aws.String("sk = :action_end"),
				TableName: &tableName,
			}
		}
	*/

	input := &dynamodb.QueryInput{
		TableName:                 &tableName,
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
	}

	result, err := dbClient.Query(input)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func main() {
	lambda.Start(handler)
}