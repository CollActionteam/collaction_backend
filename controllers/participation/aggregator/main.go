package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/CollActionteam/collaction_backend/internal/constants"
	m "github.com/CollActionteam/collaction_backend/internal/models"
	"github.com/CollActionteam/collaction_backend/models"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func groupEvents(
	events []m.ParticipationEvent,
	eventsGrouped map[string]([]m.ParticipationEvent),
	groupBy func(m.ParticipationEvent) string) {
	for _, e := range events {
		key := groupBy(e)
		if v, hasKey := eventsGrouped[key]; hasKey {
			eventsGrouped[key] = append(v, e)
		} else {
			eventsGrouped[key] = []m.ParticipationEvent{e}
		}
	}
}

func handler(ctx context.Context, sqsEvent events.SQSEvent) error {
	events := []m.ParticipationEvent{}
	for _, record := range sqsEvent.Records {
		var event m.ParticipationEvent
		fmt.Printf("Received SQS Message: %s\n", record.Body) // TODO remove!
		err := json.Unmarshal([]byte(record.Body), &event)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			events = append(events, event)
		}
	}
	eventsByCrowdaction := make(map[string]([]m.ParticipationEvent))
	groupEvents(events, eventsByCrowdaction, func(e m.ParticipationEvent) string { return e.CrowdactionID })
	for crowdactionID, crowdactionEvents := range eventsByCrowdaction {
		participantCountChangedBy := 0
		for _, event := range crowdactionEvents {
			participantCountChangedBy += event.Count
		}
		err := models.ChangeCrowdactionParticipantCountBy(crowdactionID, constants.TableName, participantCountChangedBy)
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
