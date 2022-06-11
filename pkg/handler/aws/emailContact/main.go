package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/CollActionteam/collaction_backend/internal/contact"
	"github.com/CollActionteam/collaction_backend/internal/models"
	hnd "github.com/CollActionteam/collaction_backend/pkg/handler"
	awsRepository "github.com/CollActionteam/collaction_backend/pkg/repository/aws"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
)

func handler(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	var request models.EmailContactRequest
	if err := json.Unmarshal([]byte(req.Body), &request); err != nil {
		return errToResponse(err, http.StatusBadRequest), nil
	}
	// TODO implement POW verification using nonce (see https://github.com/CollActionteam/collaction_backend/issues/58)

	if err := request.Validate(ctx); err != nil {
		body, _ := json.Marshal(hnd.Response{Status: hnd.StatusFail, Data: map[string]interface{}{"error": err}})
		return events.APIGatewayV2HTTPResponse{Body: string(body), StatusCode: http.StatusBadRequest}, nil
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

	body, _ := json.Marshal(hnd.Response{Status: hnd.StatusSuccess, Data: nil})
	return events.APIGatewayV2HTTPResponse{StatusCode: 200, Body: string(body)}, nil
}

func main() {
	lambda.Start(handler)
}

func errToResponse(err error, code int) events.APIGatewayV2HTTPResponse {
	body, _ := json.Marshal(hnd.Response{Status: hnd.StatusFail, Data: map[string]interface{}{"error": err.Error()}})
	return events.APIGatewayV2HTTPResponse{Body: string(body), StatusCode: code}
}
