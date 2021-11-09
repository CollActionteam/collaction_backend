package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/CollActionteam/collaction_backend/participation"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func groupEvents(
	events []participation.ParticipationEvent,
	eventsGrouped map[string]([]participation.ParticipationEvent),
	groupBy func(participation.ParticipationEvent) string) {
	for _, e := range events {
		key := groupBy(e)
		if v, hasKey := eventsGrouped[key]; hasKey {
			eventsGrouped[key] = append(v, e)
		} else {
			eventsGrouped[key] = []participation.ParticipationEvent{e}
		}
	}
}

func updateCollectiveCommitment(crowdactionID string, commitmentID string, changedBy int) {
	// TODO implement (Blocked by CAN-72)
	// Use atomic counter (https://docs.aws.amazon.com/amazondynamodb/latest/APIReference/API_UpdateItem.html#API_UpdateItem_Examples)
	fmt.Printf("Change participation count of %s (commitment %s) by %d\n", crowdactionID, commitmentID, changedBy)
}

func handler(ctx context.Context, kinesisEvent events.KinesisEvent) error {
	var err error

	events := []participation.ParticipationEvent{}

	for _, record := range kinesisEvent.Records {
		kinesisRecord := record.Kinesis
		dataBytes := kinesisRecord.Data
		var event participation.ParticipationEvent
		err = json.Unmarshal(dataBytes, &event)
		if err != nil {
			fmt.Print(err.Error())
		} else {
			events = append(events, event)
		}
	}

	eventsByCrowdaction := make(map[string]([]participation.ParticipationEvent))
	groupEvents(events, eventsByCrowdaction, func(e participation.ParticipationEvent) string { return e.CrowdactionID })

	for crowdactionID, crowdactionEvents := range eventsByCrowdaction {
		eventsByCommitment := make(map[string]([]participation.ParticipationEvent))
		groupEvents(crowdactionEvents, eventsByCommitment, func(e participation.ParticipationEvent) string { return e.CommitmentID })

		for commitmentID, commitmentEvents := range eventsByCommitment {
			aggregated := 0
			for _, e := range commitmentEvents {
				aggregated += e.Count
			}
			updateCollectiveCommitment(crowdactionID, commitmentID, aggregated)
		}
	}

	return nil
}

func main() {
	lambda.Start(handler)
}
