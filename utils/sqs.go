package utils

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

func CreateQueueClient() *sqs.SQS {
	sess := session.Must(session.NewSession())
	return sqs.New(sess)
}

func SendQueueMessage(qClient *sqs.SQS, queueUrl string, payload interface{}) error {
	json, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	messageBody := string(json)
	_, err = qClient.SendMessage(&sqs.SendMessageInput{
		MessageBody: &messageBody,
		QueueUrl:    aws.String(queueUrl),
	})
	return err
}
