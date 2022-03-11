package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/CollActionteam/collaction_backend/auth"
	"github.com/CollActionteam/collaction_backend/internal/constants"
	m "github.com/CollActionteam/collaction_backend/internal/models"
	"github.com/CollActionteam/collaction_backend/internal/participation"
	"github.com/CollActionteam/collaction_backend/models"
	"github.com/CollActionteam/collaction_backend/pkg/repository"
	"github.com/CollActionteam/collaction_backend/pkg/repository/aws"
	"github.com/CollActionteam/collaction_backend/utils"
	"github.com/aws/aws-lambda-go/events"
)

type ParticipationHandler struct {
	service participation.Service
}

func NewParticipationHandler() *ParticipationHandler {
	participationRepository := repository.NewParticipation(aws.NewDynamo())
	return &ParticipationHandler{service: participation.NewParticipationService(participationRepository)}
}

func (h *ParticipationHandler) createParticipation(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	userID, name, crowdaction, err := retrieveInfoFromRequest(req)
	if err != nil {
		return handlerError(err), nil
	}
	var payload m.JoinPayload
	if err := json.Unmarshal([]byte(req.Body), &payload); err != nil {
		return utils.CreateMessageHttpResponse(http.StatusInternalServerError, err.Error()), nil
	}
	if len(payload.Commitments) == 0 {
		return utils.CreateMessageHttpResponse(http.StatusBadRequest, "cannot participate without commitments"), nil
	}

	if err := models.ValidateCommitments(payload.Commitments, crowdaction.CommitmentOptions); err != nil {
		return utils.CreateMessageHttpResponse(http.StatusInternalServerError, err.Error()), nil
	}

	if err := h.service.RegisterParticipation(ctx, userID, name, crowdaction, payload); err != nil {
		return utils.CreateMessageHttpResponse(http.StatusInternalServerError, err.Error()), nil
	}
	return utils.CreateMessageHttpResponse(http.StatusOK, "updated"), nil
}

func (h *ParticipationHandler) getParticipation(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	userID, _, crowdaction, err := retrieveInfoFromRequest(req)
	if err != nil {
		return handlerError(err), nil
	}

	participation, err := h.service.GetParticipation(ctx, userID, crowdaction.CrowdactionID)
	if err != nil {
		return utils.CreateMessageHttpResponse(http.StatusInternalServerError, err.Error()), nil
	}
	if participation == nil {
		return utils.CreateMessageHttpResponse(http.StatusNotFound, "not participating"), nil
	}
	jsonPayload, _ := json.Marshal(participation)
	return events.APIGatewayV2HTTPResponse{
		Body:       string(jsonPayload),
		StatusCode: http.StatusOK,
	}, nil
}

func (h *ParticipationHandler) deleteParticipation(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	userID, _, crowdaction, err := retrieveInfoFromRequest(req)
	if err != nil {
		return handlerError(err), nil
	}
	if err := h.service.CancelParticipation(ctx, userID, crowdaction); err != nil {
		return utils.CreateMessageHttpResponse(http.StatusInternalServerError, err.Error()), nil
	}
	return utils.CreateMessageHttpResponse(http.StatusOK, "updated"), nil
}

func retrieveInfoFromRequest(req events.APIGatewayV2HTTPRequest) (string, string, *models.Crowdaction, error) {
	crowdactionID := req.PathParameters["crowdactionID"]
	usrInf, err := auth.ExtractUserInfo(req)
	if err != nil {
		return usrInf.UserID(), usrInf.Name(), nil, err
	}
	crowdaction, _ := models.GetCrowdaction(crowdactionID, constants.TableName)
	if crowdaction == nil {
		return "", "", nil, errors.New("crowdaction not found")
	}
	return usrInf.UserID(), usrInf.Name(), crowdaction, nil
}

func handlerError(err error) events.APIGatewayV2HTTPResponse {
	if err.Error() == "crowdaction not found" {
		return utils.CreateMessageHttpResponse(http.StatusNotFound, err.Error())
	}
	return utils.CreateMessageHttpResponse(http.StatusInternalServerError, err.Error())
}
