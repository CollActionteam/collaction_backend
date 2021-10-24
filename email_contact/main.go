package main

import (
	"encoding/json"
	"fmt"

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

	//just for debugging purposes
	request_details(req)

	sess := session.Must(session.NewSession())

	// Create an SES session.
	svc := ses.New(sess)

	// Assemble the email.
	input, err := buildEmail(req)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 404}, nil
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

func main() {
	lambda.Start(handler)
}

func request_details(r events.APIGatewayProxyRequest) {

	var bodyRequest Mail

	fmt.Println("[handler]...")

	fmt.Println("events.APIGatewayProxyRequest is", r)

	body := r.Body
	fmt.Println("request.Body is", body)

	fmt.Println("Headers:")
	for key, value := range r.Headers {
		fmt.Printf("    %s: %s\n", key, value)
	}

	_ = json.Unmarshal([]byte(r.Body), &bodyRequest)

	fmt.Println("bodyRequest.Email:", bodyRequest.Email)
	fmt.Println("bodyRequest.Subject:", bodyRequest.Subject)
	fmt.Println("bodyRequest.Message:", bodyRequest.Message)
	fmt.Println("bodyRequest.AppVersion", bodyRequest.AppVersion)

	fmt.Println("recipient:", getRecipient())
	fmt.Println("sender:", getSender())

}
