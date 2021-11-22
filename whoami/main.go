package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/CollActionteam/collaction_backend/auth"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(request events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
	msg := "no user info was extracted"
	statusCode := http.StatusBadRequest
	usrInf, err := auth.ExtractUserInfo(req)
	if err != nil {
		msg = err.Error()
	} else if usrInf != nil {
		msg = fmt.Sprintf("your name is %s, your user id is %s and your phone number is %s", usrInf.Name(), usrInf.UserID(), usrInf.PhoneNumber())
		statusCode = http.StatusOK
	}
	body, _ := json.Marshal(map[string]interface{}{"message": msg})
	return events.APIGatewayProxyResponse{
		Body:       string(body),
		StatusCode: statusCode,
	}, nil
}

func main() {
	lambda.Start(handler)
}
