package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/CollActionteam/collaction_backend/internal/contact"
	"github.com/CollActionteam/collaction_backend/internal/models"
	awsRepository "github.com/CollActionteam/collaction_backend/pkg/repository/aws"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/go-playground/validator/v10"
)

func handler(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	var request models.EmailContactRequest
	if err := json.Unmarshal([]byte(req.Body), &request); err != nil {
		return errToResponse(err, http.StatusBadRequest), nil
	}

	fmt.Println("Hello World")

	validate := validator.New()
	if err := validate.StructCtx(ctx, request); err != nil {
		//TODO 10.01.22 mrsoftware: fix the error message
		return errToResponse(err, http.StatusBadRequest), nil
	}

	sess := session.Must(session.NewSession())
	emailRepository := awsRepository.NewEmail(sess)
	configManager := awsRepository.NewConfigManager(sess)

	stage := req.RequestContext.Stage
	if stage == "" {
		stage = "dev"
	}

	if err := contact.NewContactService(emailRepository, configManager, stage).SendEmail(ctx, request); err != nil {
		return errToResponse(err, http.StatusInternalServerError), nil
	}

	msg, _ := json.Marshal(map[string]string{"message": "message sent successfully"})
	return events.APIGatewayV2HTTPResponse{StatusCode: 200, Body: string(msg)}, nil
}

func main() {
	lambda.Start(handler)
}

func errToResponse(err error, code int) events.APIGatewayV2HTTPResponse {
	msg, _ := json.Marshal(map[string]string{"message": err.Error()})
	return events.APIGatewayV2HTTPResponse{Body: string(msg), StatusCode: code}
}
