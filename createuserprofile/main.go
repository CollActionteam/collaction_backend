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
	var msg []byte

	fmt.Println("POST")

	err := profileservice.CreateProfile(req)
	if err != nil {

		msg, _ = json.Marshal(map[string]interface{}{"message": "Error Processing Request"})
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       string(msg),
		}, err

	}

	msg, _ = json.Marshal(map[string]interface{}{"message": "Profile Created"})
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(msg),
	}, nil
}

func main() {
	lambda.Start(handler)
}
