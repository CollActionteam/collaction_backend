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
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/go-playground/validator/v10"
)

// func handler(ctx context.Context, req events.APIGatewayV2HTTPRequest) {
func handler(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	// func handler(ctx context.Context, req events.APIGatewayV2HTTPRequest) {
	crowdactionID := req.PathParameters["crowdactionID"]

	fmt.Println("Hello Go, first message from CollAction AWS:  ", crowdactionID)
	var request models.CrowdactionRequest

	// if err := json.Unmarshal([]byte(req.Body), &request); err != nil {
	// 	return errToResponse(err, http.StatusBadRequest), nil
	// }

	fmt.Println("bad request")

	validate := validator.New()
	if err := validate.StructCtx(ctx, request); err != nil {
		body, _ := json.Marshal(hnd.Response{Status: hnd.StatusFail, Data: map[string]interface{}{"error": utils.ValidationResponse(err, validate)}})
		return events.APIGatewayV2HTTPResponse{Body: string(body), StatusCode: http.StatusBadRequest}, nil
	}

	fmt.Println("Request has been validated!")

	dynamoRepository := awsRepository.NewDynamo()

	// crowdactionService := cwd.NewCrowdactionService(dynamoRepository)

	// getCrowdaction, err := crowdactionService.GetCrowdaction(ctx, crowdactionID)

	getCrowdaction, err := cwd.NewCrowdactionService(dynamoRepository).GetCrowdaction(ctx, crowdactionID)

	if err != nil {
		// fmt.Println("There was an error")
		return utils.CreateMessageHttpResponse(http.StatusInternalServerError, err.Error()), nil
	}
	if getCrowdaction == nil {
		// fmt.Println("No crowdaction available")
		return utils.CreateMessageHttpResponse(http.StatusNotFound, "not participating"), nil
	}

	fmt.Println("getCrowdaction", getCrowdaction)

	jsonPayload, _ := json.Marshal(getCrowdaction)
	return events.APIGatewayV2HTTPResponse{
		Body:       string(jsonPayload),
		StatusCode: http.StatusOK,
	}, nil
}

func main() {
	lambda.Start(handler)
}

func errToResponse(err error, code int) events.APIGatewayV2HTTPResponse {
	msg, _ := json.Marshal(map[string]string{"message": err.Error()})
	return events.APIGatewayV2HTTPResponse{Body: string(msg), StatusCode: code}
}
