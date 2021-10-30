package main

import (
	"encoding/json"
	"net/http"

	"github.com/CollActionteam/collaction_backend/profileservice"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var msg profileservice.Response

	profileData, err := profileservice.GetProfile(req)
	if err != nil {
		msg = profileservice.Response{
			Message: "Error Retreiving Profile",
			Data:    "",
			Status:  http.StatusInternalServerError,
		}
		jsonData, _ := json.Marshal(msg)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       string(jsonData),
		}, err

	}

	if profileData == nil {
		msg = profileservice.Response{
			Message: "no user Profile found",
			Data:    "",
			Status:  404,
		}
		jsonData, _ := json.Marshal(msg)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusNotFound,
			Body:       string(jsonData),
		}, nil
	}

	msg = profileservice.Response{
		Message: "Successfully Retrieved Profile",
		Data:    profileData,
		Status:  200,
	}
	jsonData, _ := json.Marshal(msg)

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(jsonData),
	}, err
}

func main() {
	lambda.Start(handler)
}
