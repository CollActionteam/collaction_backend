package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	cwd "github.com/CollActionteam/collaction_backend/internal/crowdactions"
	"github.com/CollActionteam/collaction_backend/internal/models"
	hnd "github.com/CollActionteam/collaction_backend/pkg/handler"
	awsRepository "github.com/CollActionteam/collaction_backend/pkg/repository/aws"
	"github.com/CollActionteam/collaction_backend/utils"
	"github.com/aws/aws-lambda-go/events"
	"github.com/go-playground/validator/v10"
)

type CrowdactionHandler struct {
	service cwd.Service
}

func NewCrowdactionHandler() *CrowdactionHandler {
	crowdactionParticipation := awsRepository.NewCrowdaction(awsRepository.NewDynamo())
	return &CrowdactionHandler{service: cwd.NewCrowdactionService(crowdactionParticipation)}
}

func (c *CrowdactionHandler) createCrowdaction(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	// validate request
	// register crowdaction
	// return call to client
	return utils.CreateMessageHttpResponse(http.StatusNotImplemented, "not implemented"), nil
}

func (c *CrowdactionHandler) getCrowdaction(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	crowdactionID := req.PathParameters["crowdactionID"]
	fmt.Println("1. handler/aws/crowdaction/crowdaction.go/getCrowdaction")
	var request models.CrowdactionRequest

	validate := validator.New()
	if err := validate.StructCtx(ctx, request); err != nil {
		body, _ := json.Marshal(hnd.Response{Status: hnd.StatusFail, Data: map[string]interface{}{"error": utils.ValidationResponse(err, validate)}})
		return events.APIGatewayV2HTTPResponse{Body: string(body), StatusCode: http.StatusBadRequest}, nil
	}

	if crowdactionID == "" {
		status := req.QueryStringParameters["status"]
		return getCrowdactionsByStatus(ctx, status)
	}

	return c.getCrowdactionByID(ctx, crowdactionID)
}

func (c *CrowdactionHandler) getCrowdactionByID(ctx context.Context, crowdactionID string) (events.APIGatewayV2HTTPResponse, error) {
	// func getCrowdactionByID(ctx context.Context, crowdactionID string) (events.APIGatewayV2HTTPResponse, error) {
	// dynamoRepository := awsRepository.NewCrowdaction(awsRepository.NewDynamo())
	fmt.Println("2. handler/aws/crowdaction/crowdaction.go/getCrowdactionByID")
	getCrowdaction, err := c.service.GetCrowdactionById(ctx, crowdactionID)
	// getCrowdaction, err := cwd.NewCrowdactionService(dynamoRepository).GetCrowdactionById(ctx, crowdactionID)

	if err != nil {
		return utils.CreateMessageHttpResponse(http.StatusInternalServerError, err.Error()), nil
	}
	if getCrowdaction == nil {
		return utils.CreateMessageHttpResponse(http.StatusNotFound, "not participating"), nil
	}

	fmt.Println("7. handler/aws/crowdaction/crowdaction.go/getCrowdactionByID", getCrowdaction)

	jsonPayload, _ := json.Marshal(getCrowdaction)
	return events.APIGatewayV2HTTPResponse{
		Body:       string(jsonPayload),
		StatusCode: http.StatusOK,
	}, nil
}

func getCrowdactionsByStatus(ctx context.Context, status string) (events.APIGatewayV2HTTPResponse, error) {
	dynamoRepository := awsRepository.NewCrowdaction(awsRepository.NewDynamo())
	getCrowdactions, err := cwd.NewCrowdactionService(dynamoRepository).GetCrowdactionsByStatus(ctx, status, nil)

	if err != nil {
		return utils.CreateMessageHttpResponse(http.StatusInternalServerError, err.Error()), nil
	}
	jsonPayload, _ := json.Marshal(getCrowdactions)

	return events.APIGatewayV2HTTPResponse{
		Body:       string(jsonPayload),
		StatusCode: http.StatusOK,
	}, nil
}

func errToResponse(err error, code int) events.APIGatewayV2HTTPResponse {
	msg, _ := json.Marshal(map[string]string{"message": err.Error()})
	return events.APIGatewayV2HTTPResponse{Body: string(msg), StatusCode: code}
}
