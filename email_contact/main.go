package main

import (
	"encoding/json"
	"fmt"
	"net/mail"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

const (
	// character encoding for the email.
	charSet = "UTF-8"

	// separator between actual email message and app version
	separator = "  ### app version: "
)

// BodyRequest is our self-made struct to process JSON request from Client
type Mail struct {
	Email      string `json:"email"`
	Subject    string `json:"subject"`
	Message    string `json:"message"`
	AppVersion string `json:"app_version"`
}

func handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	sess := session.Must(session.NewSession())

	// Create an SES session.
	svc := ses.New(sess)

	// Assemble the email.
	input, err := buildEmail(req)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 400}, nil
	}

	// Attempt to send the email.
	result, err := svc.SendEmail(input)
	if err != nil {
		fmt.Println(err.Error())
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       *result.MessageId,
	}, nil
}

func buildEmail(r events.APIGatewayProxyRequest) (*ses.SendEmailInput, error) {

	var bodyRequest Mail

	// unmarshal the json, return 404 if error
	err := json.Unmarshal([]byte(r.Body), &bodyRequest)
	if err != nil {
		return nil, err
	}

	err = valid(bodyRequest.Email)
	if err != nil {
		return nil, err
	}

	return &ses.SendEmailInput{
		Destination: &ses.Destination{
			//		CcAddresses: []*string{},
			ToAddresses: []*string{
				aws.String(getRecipient()),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Text: &ses.Content{
					Charset: aws.String(charSet),
					Data: aws.String(bodyRequest.Message +
						separator +
						bodyRequest.AppVersion),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(charSet),
				Data:    aws.String(bodyRequest.Subject),
			},
		},
		Source:           aws.String(getSender()),
		ReplyToAddresses: []*string{aws.String(bodyRequest.Email)},
	}, nil

}

//email used for sender
func getSender() string {
	//TO DO obtain address from the parameter store
	return "hello@collaction.org"
}

//email used for recipient
func getRecipient() string {
	//TO DO obtain address from the parameter store
	return "hello@collaction.org"
}

func valid(email string) error {
	_, err := mail.ParseAddress(email)
	return err
}

func main() {
	lambda.Start(handler)
}
