package main

import (
	"encoding/json"

	"github.com/CollActionteam/collaction_backend/models"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kinesis"
)

func recordEvent(sess *session.Session, userID string, crowdactionID string, commitments []string, count int) error {
	kc := kinesis.New(sess)
	json, err := json.Marshal(models.ParticipationEvent{
		UserID:        userID,
		CrowdactionID: crowdactionID,
		Commitments:   commitments,
		Count:         count,
	})
	if err != nil {
		return err
	}
	_, err = kc.PutRecord(&kinesis.PutRecordInput{
		StreamName:   &streamName,
		PartitionKey: &crowdactionID,
		Data:         json,
	})
	return err
}
