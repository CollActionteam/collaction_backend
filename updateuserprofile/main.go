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

	err := profileservice.UpdateProfile(req)

	if err != nil {

		msg, _ = json.Marshal(map[string]interface{}{"message": "Error Processing Update Request"})
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       string(msg),
		}, err

	}

	msg, _ = json.Marshal(map[string]interface{}{"message": "profile update successful"})
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(msg),
	}, nil
}

func main() {
	lambda.Start(handler)
}
