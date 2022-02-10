package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/CollActionteam/collaction_backend/auth"
	"github.com/CollActionteam/collaction_backend/internal/uploads"
	awsRepository "github.com/CollActionteam/collaction_backend/pkg/repository/aws"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
)

func handler(ctx context.Context, req events.APIGatewayV2HTTPRequest) (res events.APIGatewayProxyResponse, err error) {

	// This is to recover from a "invalid memory address or nil pointer dereference: errorString" runtime error
	// when the endpoint is called without valid auth token
	defer func() {
		if r := recover(); r != nil {
			res = events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}
			err = fmt.Errorf(fmt.Sprintf("%v", r))
			return
		}
	}()

	usrInf, err := auth.ExtractUserInfo(req)
	if err != nil {
		res = events.APIGatewayProxyResponse{StatusCode: http.StatusForbidden, Body: "user not authorized"}
		return res, err
	}

	// getting user id, which will be used as object key
	userID := usrInf.UserID()
	sess := session.Must(session.NewSession())
	profileImageUploadRepo := awsRepository.NewProfilePicture(sess)

	strUrl, err := uploads.NewProfileImageUploadService(profileImageUploadRepo).GetUploadUrl("png", userID)
	if err != nil {
		res = events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError, Body: "error generating link"}
		return res, err
	}

	response, _ := json.Marshal(map[string]interface{}{"upload_url": strUrl})

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(response),
	}, nil
}

func main() {
	lambda.Start(handler)
}
