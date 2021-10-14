package main

import (
	"encoding/json"
	"fmt"

	"github.com/CollActionteam/collaction_backend/auth"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	msg := "No user info was extracted!"
	statusCode := 400
	u := auth.ExtractUserInfo(request)
	if u != nil {
		msg = fmt.Sprintf("You name is %s, your user id is %s and your phone number is %s", u.Name, u.UserID, u.PhoneNumber)
		statusCode = 200
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
