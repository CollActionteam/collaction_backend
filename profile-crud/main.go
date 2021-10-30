package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/CollActionteam/collaction_backend/profileservice"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func getProfile(req events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
	var msg profileservice.Response

	profileData, err := profileservice.GetProfile(req)
	if err != nil {
		msg = profileservice.Response{
			Message: "Error Retreiving Profile",
			Data:    "",
			Status:  http.StatusInternalServerError,
		}
		jsonData, _ := json.Marshal(msg)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       string(jsonData),
		}, err

	}

	if profileData == nil {
		msg = profileservice.Response{
			Message: "no user Profile found",
			Data:    "",
			Status:  404,
		}
		jsonData, _ := json.Marshal(msg)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusNotFound,
			Body:       string(jsonData),
		}, nil
	}

	msg = profileservice.Response{
		Message: "Successfully Retrieved Profile",
		Data:    profileData,
		Status:  200,
	}
	jsonData, _ := json.Marshal(msg)

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(jsonData),
	}, err
}

func createProfile(req events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
	err := profileservice.CreateProfile(req)
	if err != nil {
		tmsg := profileservice.Response{
			Message: fmt.Sprintf("%v", err),
			Data:    "",
			Status:  http.StatusBadRequest,
		}
		jsonData, _ := json.Marshal(tmsg)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       string(jsonData),
		}, nil

	}

	tmsg := profileservice.Response{
		Message: "Profile Created",
		Data:    "",
		Status:  200,
	}
	jsonData, _ := json.Marshal(tmsg)
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(jsonData),
	}, nil
}

func updateProfile(req events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
	err := profileservice.UpdateProfile(req)

	if err != nil {
		tmsg := profileservice.Response{
			Message: fmt.Sprintf("%v", err),
			Data:    "",
			Status:  200,
		}
		jsonData, _ := json.Marshal(tmsg)

		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       string(jsonData),
		}, err

	}

	tmsg := profileservice.Response{
		Message: "profile update successful",
		Data:    "",
		Status:  http.StatusBadRequest,
	}
	jsonData, _ := json.Marshal(tmsg)
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(jsonData),
	}, nil
}

func handler(req events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
	method := strings.ToLower(req.RequestContext.HTTP.Method)

	var res events.APIGatewayProxyResponse
	var err error

	if method == "get" {
		res, err = getProfile(req)
	} else if method == "post" {
		res, err = createProfile(req)
	} else if method == "put" {
		res, err = updateProfile(req)
	} else {
		jsonData, _ := json.Marshal(map[string]interface{}{"message": "Not implemented"})
		res = events.APIGatewayProxyResponse{
			StatusCode: 501,
			Body:       string(jsonData),
		}
	}

	return res, err
}

func main() {
	lambda.Start(handler)
}
