package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/CollActionteam/collaction_backend/auth"
	pps "github.com/CollActionteam/collaction_backend/profilepictureservice"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (res events.APIGatewayProxyResponse, err error) {
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

	strUrl, err := pps.GetUploadUrl(req.PathParameters["ext"], userID)
	if err != nil {
		res = events.APIGatewayProxyResponse{StatusCode: http.StatusForbidden, Body: "user not authorized"}
		return res, err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       strUrl,
	}, nil
}

func main() {
	lambda.Start(handler)
}
