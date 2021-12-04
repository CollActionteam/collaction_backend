package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/CollActionteam/collaction_backend/models"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	tableName = os.Getenv("TABLE_NAME")
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

func handler(ctx context.Context, sqsEvent events.SQSEvent) error {
	events := []models.ParticipationEvent{}
	for _, record := range sqsEvent.Records {
		var event models.ParticipationEvent
		fmt.Printf("Received SQS Message: %s\n", record.Body) // TODO remove!
		err := json.Unmarshal([]byte(record.Body), &event)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			events = append(events, event)
		}
	}
	eventsByCrowdaction := make(map[string]([]models.ParticipationEvent))
	groupEvents(events, eventsByCrowdaction, func(e models.ParticipationEvent) string { return e.CrowdactionID })
	for crowdactionID, crowdactionEvents := range eventsByCrowdaction {
		participantCountChangedBy := 0
		for _, event := range crowdactionEvents {
			participantCountChangedBy += event.Count
		}
		err := models.ChangeCrowdactionParticipantCountBy(crowdactionID, tableName, participantCountChangedBy)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Printf("Change participation count of %s by %d\n", crowdactionID, participantCountChangedBy)
		}
		// TODO count individual commitments and store them in the crowdaction record
	}
	return nil
}

func main() {
	lambda.Start(handler)
}
