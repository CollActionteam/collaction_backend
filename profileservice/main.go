package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/collactionteam/collaction_backend/handlers"
)

func main() {
	lambda.Start(handlers.CreateProfileHandler)
}
