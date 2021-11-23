package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/CollActionteam/collaction_backend/auth"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func getUploadUrl(ext string, userID string) (string, error) {

	var (
		bucket  = os.Getenv("BUCKET")
		filekey = userID + "." + ext
	)

	// Initialize a session that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials.

	// Create S3 service client
	svc := s3.New(session.Must(session.NewSession()))
	reqs, _ := svc.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filekey),
	})

	str, err := reqs.Presign(15 * time.Minute)

	if err != nil {
		return "", err
	}
	return str, nil
}

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

	strUrl, err := getUploadUrl("png", userID)
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
