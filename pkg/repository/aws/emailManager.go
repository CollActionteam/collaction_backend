package aws

import (
	"context"
	"github.com/CollActionteam/collaction_backend/internal/constants"
	"github.com/CollActionteam/collaction_backend/internal/models"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/pkg/errors"
)

type Email struct {
	Client *ses.SES
}

func NewEmail(session *session.Session) *Email {
	return &Email{Client: ses.New(session)}
}

func (e *Email) Send(ctx context.Context, data models.EmailData) error {
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{ToAddresses: []*string{aws.String(data.Recipient)}},
		Message: &ses.Message{
			Body: &ses.Body{
				Text: &ses.Content{Charset: aws.String(constants.CharSet), Data: aws.String(data.Message)},
			},
			Subject: &ses.Content{Charset: aws.String(constants.CharSet), Data: aws.String(data.Subject)},
		},
		Source:           aws.String(data.Sender),
		ReplyToAddresses: []*string{aws.String(data.ReplyEmail)},
	}

	if _, err := e.Client.SendEmail(input); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
