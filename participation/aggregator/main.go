package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/CollActionteam/collaction_backend/participation"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, kinesisEvent events.KinesisEvent) error {
	var err error
	aggregations := make(map[string]int)

	for _, record := range kinesisEvent.Records {
		kinesisRecord := record.Kinesis
		dataBytes := kinesisRecord.Data
		var e participation.ParticipationEvent
		err = json.Unmarshal(dataBytes, &e)
		if err != nil {
			fmt.Print(err.Error())
		} else {
			if v, hasKey := aggregations[e.CrowdactionID]; hasKey {
				aggregations[e.CrowdactionID] = v + e.Count
			} else {
				aggregations[e.CrowdactionID] = e.Count
			}
		}
	}

	for k, v := range aggregations {
		// TODO implement
		fmt.Printf("Change participation count of %s by %d\n", k, v)
	}

	return nil
}

func main() {
	lambda.Start(handler)
}
