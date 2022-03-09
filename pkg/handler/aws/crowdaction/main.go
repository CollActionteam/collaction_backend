package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	crowdaction "github.com/CollActionteam/collaction_backend/internal/crowdactions"
	"github.com/CollActionteam/collaction_backend/internal/models"
	hnd "github.com/CollActionteam/collaction_backend/pkg/handler"
	awsRepository "github.com/CollActionteam/collaction_backend/pkg/repository/aws"
	"github.com/CollActionteam/collaction_backend/utils"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	// "github.com/CollActionteam/collaction_backend/internal/contact"
	// "github.com/CollActionteam/collaction_backend/internal/models"
	"github.com/go-playground/validator/v10"
)

// type CrowdactionHandler struct {
// 	// getting all the main functions from the internal
// 	crowadction crowdactions.Dynamo
// }

// // return crowdactions by
// func getCrowdactions (status string) (events.APIGatewayProxyResponse, error) {}

// // return single crowdaction by ID
// func getCrowdaction (crowdactionId string, dynamoRepository Dynamo) (events.APIGatewayProxyResponse, error) {
// 	crowdaction, err := models.

// 	// should call the internal and receive a response
// }

// func retrieveInfoFromRequest(req events.APIGatewayV2HTTPRequest) (*models.Crowdaction, error) {
// 	crowdactionID := req.PathParameters["crowdactionID"]
// 	usrInf, err := auth.ExtractUserInfo(req)
// 	if err != nil {
// 		return nil, err
// 	}
// 	crowdaction, _ := models.GetCrowdaction(crowdactionID, constants.TableName)
// 	if crowdaction == nil {
// 		return "", "", nil, errors.New("crowdaction not found")
// 	}
// 	return crowdaction, nil
// }

func handler(ctx context.Context, req events.APIGatewayV2HTTPRequest) string {
	var request models.CrowdactionRequest

	if err := json.Unmarshal([]byte(req.Body), &request); err != nil {
		return errToResponse(err, http.StatusBadRequest), nil
	}

	fmt.Println("Hello Go, first message from AWS!")

	validate := validator.New()
	if err := validate.StructCtx(ctx, request); err != nil {
		body, _ := json.Marshal(hnd.Response{Status: hnd.StatusFail, Data: map[string]interface{}{"error": utils.ValidationResponse(err, validate)}})
		return events.APIGatewayV2HTTPResponse{Body: string(body), StatusCode: http.StatusBadRequest}, nil
	}

	crowdactionID := req.PathParameters["crowdactionID"]
	dynamoRepository := awsRepository.NewDynamo()

	if err := crowdaction.NewCrowdactionService(dynamoRepository).GetCrowdaction(crowdactionID)
	// fmt.Println("Getting the crowdactionID from the request")

	// if crowdactionID == "" {
	// 	status := req.QueryStringParameters["status"]
	// 	switch status {
	// 	case "":
	// 		status = "joinable"
	// 	case "featured":
	// 		status = "joinable"
	// 	}
	// 	return getCrowdactions(status)
	// }

	// return getCrowdaction(crowdactionID)
	
	return event.APIGatewayProxyResponse{StatusCode: 200, Body: string(body)}, nill
}

func main() {
	lambda.Start(handler)
}

func errToResponse(err error, code int) events.APIGatewayV2HTTPResponse {
	msg, _ := json.Marshal(map[string]string{"message": err.Error()})
	return events.APIGatewayV2HTTPResponse{Body: string(msg), StatusCode: code}
}
