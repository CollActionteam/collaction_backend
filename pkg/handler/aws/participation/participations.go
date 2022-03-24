package main

import (
	"context"
	"encoding/json"
	"net/http"

	m "github.com/CollActionteam/collaction_backend/internal/models"
	"github.com/CollActionteam/collaction_backend/utils"
	"github.com/aws/aws-lambda-go/events"
)

type participations *[]m.ParticipationRecord

func (h *ParticipationHandler) getParticipations(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	var err error
	userID, _, _, err := retrieveInfoFromRequest(req)
	if err != nil {
		return handlerError(err), nil
	}

	particpations, err := h.service.GetParticipations(ctx, userID)
	if err != nil {
		return utils.CreateMessageHttpResponse(http.StatusNotFound, "Invalid UserId"), err
	}

	jsonPayload, _ := json.Marshal(particpations)
	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusOK,
		Body:       string(jsonPayload),
	}, nil
}
