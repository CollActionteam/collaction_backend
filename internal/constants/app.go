package constants

import "os"

const (
	CharSet                  = "UTF-8"
	RecipientEmail           = "/collaction/%s/contact/email" // stage
	DisplayNameMinimumLength = 2
	DisplayNameMaximumLength = 20
	CountryMinimumLength     = 3
	CountryMaximumLength     = 20
	CityMinimumLength        = 3
	CityMaximumLength        = 20
	BioMinimumLength         = 10
	BioMaximumLength         = 100
)

var (
	TableName              = os.Getenv("TABLE_NAME")
	IndexName              = os.Getenv("INDEX_NAME")
	ProfileTablename       = os.Getenv("PROFILE_TABLE")
	ParticipationQueueName = os.Getenv("PARTICIPATION_QUEUE")
)
