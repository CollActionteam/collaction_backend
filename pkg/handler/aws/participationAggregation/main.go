package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, sqsEvent events.SQSEvent) error {
	return NewParticipationAggregationHandler().aggregateParticipations(ctx, sqsEvent)
}

func main() {
	lambda.Start(handler)
}
