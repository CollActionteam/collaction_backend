package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/CollActionteam/collaction_backend/internal/models"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/go-playground/validator/v10"
)

var (
	tableName = os.Getenv("TABLE_NAME")
)

func handler(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	fmt.Println("Entering the crowdaction endpoint")

	// invoking crowdaction model
	var request models.CrowdActionRequest

	// creating a new validator object
	validate := validator.New()

	// validating against error
	if err := validate.StructCtx(ctx, request); err != nil {
		return errToResponse(err, http.StatusBadRequest), nil
	}

	// AWS session
	//sess := session.Must(session.NewSession())

	if err := json.Unmarshal([]byte(req.Body), &request); err != nil {
		return errToResponse(err, http.StatusBadRequest), nil
	}

	return events.APIGatewayV2HTTPResponse{StatusCode: 200, Body: "Success!"}, nil
}

func main() {
	lambda.Start(handler)
}

func errToResponse(err error, code int) events.APIGatewayV2HTTPResponse {
	msg, _ := json.Marshal(map[string]string{"message": err.Error()})
	return events.APIGatewayV2HTTPResponse{Body: string(msg), StatusCode: code}
}
