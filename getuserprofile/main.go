package main

import (
	"encoding/json"
	"net/http"

	"github.com/CollActionteam/collaction_backend/profileservice"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var msg []byte

	profileData, err := profileservice.GetProfile(req)
	if err != nil {

		msg, _ = json.Marshal(map[string]interface{}{"message": "Error Retreiving Profile"})
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       string(msg),
		}, err

	}

	jsonData, err := json.Marshal(profileData)
	if err != nil {

		msg, _ = json.Marshal(map[string]interface{}{"message": "Error Encoding Data"})
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       string(msg),
		}, err

	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(jsonData),
	}, err
}

func main() {
	lambda.Start(handler)
}
