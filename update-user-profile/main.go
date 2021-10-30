package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/CollActionteam/collaction_backend/profileservice"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	err := profileservice.UpdateProfile(req)

	if err != nil {
		tmsg := profileservice.Response{
			Message: fmt.Sprintf("%v", err),
			Data:    "",
			Status:  200,
		}
		jsonData, _ := json.Marshal(tmsg)

		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       string(jsonData),
		}, err

	}

	tmsg := profileservice.Response{
		Message: "profile update successful",
		Data:    "",
		Status:  http.StatusBadRequest,
	}
	jsonData, _ := json.Marshal(tmsg)
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(jsonData),
	}, nil
}

func main() {
	lambda.Start(handler)
}
