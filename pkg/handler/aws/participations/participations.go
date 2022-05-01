package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	m "github.com/CollActionteam/collaction_backend/internal/models"
	"github.com/CollActionteam/collaction_backend/internal/participation"
	"github.com/CollActionteam/collaction_backend/pkg/repository/aws"
	"github.com/CollActionteam/collaction_backend/utils"
	"github.com/aws/aws-lambda-go/events"
)

type participations *[]m.ParticipationRecord

type ParticipationsHandler struct {
	service participation.Service
}

func NewParticipationsHandler() *ParticipationsHandler {
	participationRepository := aws.NewParticipation(aws.NewDynamo())
	return &ParticipationsHandler{service: participation.NewParticipationService(participationRepository)}
}

func (h *ParticipationsHandler) getParticipations(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	var err error
	var data *[]m.ParticipationRecord
	path := req.RequestContext.HTTP.Path
	// TODO pagination
	if strings.HasPrefix(path, "/crowdactions") {
		crowdactionID := req.PathParameters["crowdactionID"]
		// TODO check password:
		// 1. Fetch crowdaction
		// 2. If the crowdaction is not found, return 404
		// 3. If the crowdaction is password protected, check the request for the password
		data, err = h.service.GetParticipationsCrowdaction(ctx, crowdactionID)
	} else if strings.HasPrefix(path, "/profiles") {
		userID := req.PathParameters["userID"]
		data, err = h.service.GetParticipationsUser(ctx, userID)
	} else {
		err = fmt.Errorf("invalid path: %s", path)
	}
	if err != nil {
		handlerError(err)
	}
	return utils.GetDataHttpResponse(http.StatusOK, "", data), nil
}

func handlerError(err error) events.APIGatewayV2HTTPResponse {
	return utils.CreateMessageHttpResponse(http.StatusInternalServerError, err.Error())
}
