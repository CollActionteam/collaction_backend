package main

import (
	"context"
	"net/http"
	"strings"

	"github.com/CollActionteam/collaction_backend/utils"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	method := strings.ToLower(req.RequestContext.HTTP.Method)

	switch method {
	case "post":
		return NewCrowdactionHandler().createCrowdaction(ctx, req)
	case "get":
		return NewCrowdactionHandler().getCrowdaction(ctx, req)
	default:
		return utils.CreateMessageHttpResponse(http.StatusNotImplemented, "not implemented"), nil
	}
}

func main() {
	lambda.Start(handler)
}
