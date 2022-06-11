package main

import (
	"context"
	"encoding/json"
	"net/http"

	cwd "github.com/CollActionteam/collaction_backend/internal/crowdactions"
	m "github.com/CollActionteam/collaction_backend/internal/models"
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
	var payload m.CrowdactionData
	if err := json.Unmarshal([]byte(req.Body), &payload); err != nil {
		return utils.CreateMessageHttpResponse(http.StatusInternalServerError, err.Error()), nil
	}

	res, err := c.service.RegisterCrowdaction(ctx, payload)

	if err != nil {
		return utils.CreateMessageHttpResponse(http.StatusInternalServerError, err.Error()), nil

	}
	// return call to client
	body, _ := json.Marshal(hnd.Response{Status: hnd.StatusSuccess, Data: res})

	return events.APIGatewayV2HTTPResponse{
		Body:       string(body),
		StatusCode: http.StatusOK,
	}, nil
}

/**
	Gateway for all the GET crowdaction methods
**/

func (c *CrowdactionHandler) getCrowdaction(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	crowdactionID := req.PathParameters["crowdactionID"]
	var request m.CrowdactionRequest

	validate := validator.New()
	if err := validate.StructCtx(ctx, request); err != nil {
		body, _ := json.Marshal(hnd.Response{Status: hnd.StatusFail, Data: map[string]interface{}{"error": utils.ValidationResponse(err, validate)}})
		return events.APIGatewayV2HTTPResponse{Body: string(body), StatusCode: http.StatusBadRequest}, nil
	}

	if crowdactionID == "" && req.QueryStringParameters["status"] != "" { // get only by status
		status := req.QueryStringParameters["status"]
		return c.getCrowdactionsByStatus(ctx, status)
	} else if crowdactionID == "" { // to get all crowdactions
		return c.getAllCrowdactions(ctx)
	}

	return c.getCrowdactionByID(ctx, crowdactionID) // get crowdaction by id
}

func (c *CrowdactionHandler) getAllCrowdactions(ctx context.Context) (events.APIGatewayV2HTTPResponse, error) {
	getCrowdactions, err := c.service.GetAllCrowdactions(ctx)
	if err != nil {
		return utils.CreateMessageHttpResponse(http.StatusInternalServerError, err.Error()), nil
	}

	jsonPayload, _ := json.Marshal(getCrowdactions)

	return events.APIGatewayV2HTTPResponse{
		Body:       string(jsonPayload),
		StatusCode: http.StatusOK,
	}, nil
}

func (c *CrowdactionHandler) getCrowdactionByID(ctx context.Context, crowdactionID string) (events.APIGatewayV2HTTPResponse, error) {
	getCrowdaction, err := c.service.GetCrowdactionById(ctx, crowdactionID)

	if err != nil {
		return utils.CreateMessageHttpResponse(http.StatusInternalServerError, err.Error()), nil
	}
	if getCrowdaction == nil {
		return utils.CreateMessageHttpResponse(http.StatusNotFound, "Crowdaction not found!"), nil
	}

	jsonPayload, _ := json.Marshal(getCrowdaction)
	return events.APIGatewayV2HTTPResponse{
		Body:       string(jsonPayload),
		StatusCode: http.StatusOK,
	}, nil
}

func (c *CrowdactionHandler) getCrowdactionsByStatus(ctx context.Context, status string) (events.APIGatewayV2HTTPResponse, error) {
	getCrowdactions, err := c.service.GetCrowdactionsByStatus(ctx, status, nil)

	if err != nil {
		return utils.CreateMessageHttpResponse(http.StatusInternalServerError, err.Error()), nil
	}
	jsonPayload, _ := json.Marshal(getCrowdactions)

	return events.APIGatewayV2HTTPResponse{
		Body:       string(jsonPayload),
		StatusCode: http.StatusOK,
	}, nil
}
