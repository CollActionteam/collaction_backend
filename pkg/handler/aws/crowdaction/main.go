package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	// "github.com/CollActionteam/collaction_backend/internal/contact"
	// "github.com/CollActionteam/collaction_backend/internal/models"
)

func handler(ctx context.Context, req events.APIGatewayV2HTTPRequest) string {
	fmt.Println("Hello Go, first message from AWS!")

	// var request models.CrowdActionRequest

	crowdactionID := req.PathParameters["crowdactionID"]

	fmt.Println("Getting the crowdactionID from the request")

	if crowdactionID == "" {
		firstReturnValue := "Returning all crowdactions"
		return firstReturnValue
	}
	firstReturnValue := "Returning specific crowdaction"

	return firstReturnValue
}

func main() {
	lambda.Start(handler)
}

func errToResponse(err error, code int) events.APIGatewayV2HTTPResponse {
	msg, _ := json.Marshal(map[string]string{"message": err.Error()})
	return events.APIGatewayV2HTTPResponse{Body: string(msg), StatusCode: code}
}
