package main

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	msg := "No jwt claims!"
	statusCode := 400
	if claims, hasClaims := request.RequestContext.Authorizer["jwt"]; hasClaims {
		jsonClaims, _ := json.Marshal(claims)
		msg = fmt.Sprintf("Extracted claims: %s", string(jsonClaims))
		statusCode = 200
	}
	body, _ := json.Marshal(map[string]interface{}{"message": msg})
	return events.APIGatewayProxyResponse{
		Body:       string(body),
		StatusCode: statusCode,
	}, nil
}

func main() {
	lambda.Start(handler)
}
