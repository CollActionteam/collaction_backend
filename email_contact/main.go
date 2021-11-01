package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/mail"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/aws/aws-sdk-go/service/ssm"
)

const (
	// character encoding for the email.
	charSet = "UTF-8"

	// separator between actual email message and app version
	separator          = "  ### app version: "
	max_subject_length = 50
	max_message_length = 500
)

// BodyRequest is our self-made struct to process JSON request from Client
type Mail struct {
	Email      string `json:"email"`
	Subject    string `json:"subject"`
	Message    string `json:"message"`
	AppVersion string `json:"app_version"`
}

var sess *session.Session

func handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	sess = session.Must(session.NewSession())

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

	//validate input
	err = isValid("email", bodyRequest.Email)
	if err != nil {
		return nil, err
	}

	err = isValid("subject", bodyRequest.Subject)
	if err != nil {
		return nil, err
	}

	err = isValid("message", bodyRequest.Message)
	if err != nil {
		return nil, err
	}

	stage := r.RequestContext.Stage
	if stage == "" {
		stage = "dev"
	}
	paramName := "/collaction/" + stage + "/contact/email"
	recipient, err := getParameterValue(paramName)
	if err != nil {
		return nil, err
	}
	if recipient == "" {
		return nil, errors.New("no email value")
	}
	sender := recipient

	return &ses.SendEmailInput{
		Destination: &ses.Destination{
			//		CcAddresses: []*string{},
			ToAddresses: []*string{
				aws.String(recipient),
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
		Source:           aws.String(sender),
		ReplyToAddresses: []*string{aws.String(bodyRequest.Email)},
	}, nil

}

func isValid(input string, value string) error {

	switch input {
	case "email":
		_, err := mail.ParseAddress(value)
		return err
	case "subject":
		if len(value) > max_subject_length {
			return errors.New("email subject is more than " + fmt.Sprint(max_subject_length) + " characters")
		}
	case "mesage":
		if len(value) > max_message_length {
			return errors.New("email message is more than " + fmt.Sprint(max_message_length) + " characters")
		}
	}

	return nil

}

func getParameterValue(paramName string) (string, error) {

	ssmsvc := ssm.New(sess)
	param, err := ssmsvc.GetParameter(&ssm.GetParameterInput{
		Name: aws.String(paramName),
	})
	if err != nil {
		return "", err
	}

	return *param.Parameter.Value, nil
}

func main() {
	lambda.Start(handler)
}
