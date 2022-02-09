package main

import (
	"context"
	"encoding/json"
	"github.com/CollActionteam/collaction_backend/auth"
	"github.com/CollActionteam/collaction_backend/internal/constants"
	m "github.com/CollActionteam/collaction_backend/internal/models"
	"github.com/CollActionteam/collaction_backend/internal/participation"
	"github.com/CollActionteam/collaction_backend/models"
	"github.com/CollActionteam/collaction_backend/pkg/repository"
	"github.com/CollActionteam/collaction_backend/pkg/repository/aws"
	"github.com/CollActionteam/collaction_backend/utils"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"net/http"
	"strings"
)

func handler(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
	method := strings.ToLower(req.RequestContext.HTTP.Method)
	crowdactionID := req.PathParameters["crowdactionID"]
	var crowdaction *models.Crowdaction

	usrInf, err := auth.ExtractUserInfo(req)
	if err != nil {
		return utils.GetMessageHttpResponse(http.StatusInternalServerError, err.Error()), nil
	}

	crowdaction, _ = models.GetCrowdaction(crowdactionID, constants.TableName)
	if crowdaction == nil {
		return utils.GetMessageHttpResponse(http.StatusNotFound, "crowdaction not found"), nil
	}

	participationRepository := repository.NewParticipation(aws.NewDynamo())
	participationService := participation.NewParticipationService(participationRepository)

	switch method {
	case "post":
		var payload m.JoinPayload
		err = json.Unmarshal([]byte(req.Body), &payload)
		if err == nil {
			if len(payload.Commitments) == 0 {
				return utils.GetMessageHttpResponse(http.StatusBadRequest, "cannot participate without commitments"), nil
			}
			err = participationService.RegisterParticipation(ctx, usrInf.UserID(), usrInf.Name(), crowdaction, payload)
		}
	case "delete":
		err = participationService.CancelParticipation(ctx, usrInf.UserID(), crowdaction)
	case "get":
		participation, err := participationService.GetParticipation(ctx, usrInf.UserID(), crowdactionID)
		if err != nil {
			return utils.GetMessageHttpResponse(http.StatusInternalServerError, err.Error()), nil
		}
		var res events.APIGatewayProxyResponse
		if participation == nil {
			res = utils.GetMessageHttpResponse(http.StatusNotFound, "not participating")
		} else {
			// "Cannot go wrong"
			jsonPayload, _ := json.Marshal(participation)
			res = events.APIGatewayProxyResponse{
				Body:       string(jsonPayload),
				StatusCode: http.StatusOK,
			}
		}
		return res, nil
	default:
		return utils.GetMessageHttpResponse(http.StatusNotImplemented, "not implemented"), nil

	}

	if err != nil {
		return utils.GetMessageHttpResponse(http.StatusInternalServerError, err.Error()), nil
	} else {
		return utils.GetMessageHttpResponse(http.StatusOK, "updated"), nil
	}

}

func main() {
	lambda.Start(handler)
}
