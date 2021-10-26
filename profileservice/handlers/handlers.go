package handlers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/collactionteam/collaction_backend/dynamodb"
	"github.com/collactionteam/collaction_backend/models"
)

func CreateProfileHandler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var profiledata models.Profile

	err := json.Unmarshal([]byte(req.Body), &profiledata)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       "Invalid Profile",
		}, nil
	}

	err = dynamodb.InsertItemIntoTable(profiledata, "profile")
	if err != nil {
		fmt.Println("Unable to insert profile into tabele")
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
	}, nil
}
