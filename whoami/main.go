package main

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	msg := "No jwt claims!"
	statusCode := 400
	var claims interface{} = nil
	if authJWT, hasAuthJWT := request.RequestContext.Authorizer["jwt"]; hasAuthJWT {
		authJWTAsMap, _ := authJWT.(map[string]interface{})
		claims = authJWTAsMap["claims"]
		msg = "Extracted claims!"
		statusCode = 200
	}
	body, _ := json.Marshal(map[string]interface{}{"message": msg, "claims": claims})
	return events.APIGatewayProxyResponse{
		Body:       string(body),
		StatusCode: statusCode,
	}, nil
}

func main() {
	lambda.Start(handler)
}
