package constants

import "os"

const (
	CharSet        = "UTF-8"
	RecipientEmail = "/collaction/%s/contact/email" // stage

)

var (
	TableName              = os.Getenv("TABLE_NAME")
	ParticipationQueueName = os.Getenv("PARTICIPATION_QUEUE")
)
