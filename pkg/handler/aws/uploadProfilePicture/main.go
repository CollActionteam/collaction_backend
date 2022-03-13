package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/CollActionteam/collaction_backend/auth"
	"github.com/CollActionteam/collaction_backend/internal/uploads"
	awsRepository "github.com/CollActionteam/collaction_backend/pkg/repository/aws"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
)

func handler(ctx context.Context, req events.APIGatewayV2HTTPRequest) (res events.APIGatewayV2HTTPResponse, err error) {
	usrInf, err := auth.ExtractUserInfo(req)
	if err != nil {
		res = events.APIGatewayV2HTTPResponse{StatusCode: http.StatusForbidden, Body: "user not authorized"}
		return res, err
	}

	userID := usrInf.UserID()
	sess := session.Must(session.NewSession())
	profileImageUploadRepo := awsRepository.NewProfilePicture(sess)

	strUrl, err := uploads.NewProfileImageUploadService(profileImageUploadRepo).GetUploadUrl(ctx, "png", userID)
	if err != nil {
		res = events.APIGatewayV2HTTPResponse{StatusCode: http.StatusInternalServerError, Body: "error generating link"}
		return res, err
	}

	response, _ := json.Marshal(map[string]interface{}{"upload_url": strUrl})

	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusOK,
		Body:       string(response),
	}, nil
}

func main() {
	lambda.Start(handler)
}
