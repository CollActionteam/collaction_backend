package main

import (
	"context"
	"github.com/CollActionteam/collaction_backend/utils"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"net/http"
	"strings"
)

func handler(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	method := strings.ToLower(req.RequestContext.HTTP.Method)
	switch method {
	case "post":
		return NewContactHandler().createParticipation(ctx, req)
	case "delete":
		return NewContactHandler().deleteParticipation(ctx, req)
	case "get":
		return NewContactHandler().getParticipation(ctx, req)
	default:
		return utils.CreateMessageHttpResponse(http.StatusNotImplemented, "not implemented"), nil

	}
}

func main() {
	lambda.Start(handler)
}
