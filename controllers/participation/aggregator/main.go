package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/CollActionteam/collaction_backend/models"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func groupEvents(
	events []models.ParticipationEvent,
	eventsGrouped map[string]([]models.ParticipationEvent),
	groupBy func(models.ParticipationEvent) string) {
	for _, e := range events {
		key := groupBy(e)
		if v, hasKey := eventsGrouped[key]; hasKey {
			eventsGrouped[key] = append(v, e)
		} else {
			eventsGrouped[key] = []models.ParticipationEvent{e}
		}
	}
}

func updateCollectiveCommitment(crowdactionID string, commitmentID string, changedBy int) {
	// TODO Use atomic counter (https://docs.aws.amazon.com/amazondynamodb/latest/APIReference/API_UpdateItem.html#API_UpdateItem_Examples)
	fmt.Printf("Change participation count of %s (commitment %s) by %d\n", crowdactionID, commitmentID, changedBy)
}

func handler(ctx context.Context, kinesisEvent events.KinesisEvent) error {
	var err error

	events := []models.ParticipationEvent{}

	for _, record := range kinesisEvent.Records {
		kinesisRecord := record.Kinesis
		dataBytes := kinesisRecord.Data
		var event models.ParticipationEvent
		err = json.Unmarshal(dataBytes, &event)
		if err != nil {
			fmt.Print(err.Error())
		} else {
			events = append(events, event)
		}
	}

	eventsByCrowdaction := make(map[string]([]models.ParticipationEvent))
	groupEvents(events, eventsByCrowdaction, func(e models.ParticipationEvent) string { return e.CrowdactionID })

	/* TODO update code for new ParticipationEvent type
	for crowdactionID, crowdactionEvents := range eventsByCrowdaction {
		eventsByCommitment := make(map[string]([]models.ParticipationEvent))
		groupEvents(crowdactionEvents, eventsByCommitment, func(e models.ParticipationEvent) string { return e.CommitmentID })

		for commitmentID, commitmentEvents := range eventsByCommitment {
			aggregated := 0
			for _, e := range commitmentEvents {
				aggregated += e.Count
			}
			updateCollectiveCommitment(crowdactionID, commitmentID, aggregated)
		}
	}
	*/

	return nil
}

func main() {
	lambda.Start(handler)
}
