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
			Status:  http.StatusNotFound,
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
		Status:  http.StatusOK,
	}
	jsonData, _ := json.Marshal(msg)

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(jsonData),
	}, err
}

func getProfileByID(req events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
	var msg profileservice.Response

	profileData, err := profileservice.GetProfileByID(req)
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
			Status:  http.StatusNotFound,
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
		Status:  http.StatusOK,
	}
	jsonData, _ := json.Marshal(msg)

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
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
		Status:  http.StatusCreated,
	}
	jsonData, _ := json.Marshal(tmsg)
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusCreated,
		Body:       string(jsonData),
	}, nil
}

func updateProfile(req events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
	err := profileservice.UpdateProfile(req)

	if err != nil {
		tmsg := profileservice.Response{
			Message: fmt.Sprintf("%v", err),
			Data:    "",
			Status:  http.StatusInternalServerError,
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
		Status:  http.StatusOK,
	}
	jsonData, _ := json.Marshal(tmsg)
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(jsonData),
	}, nil
}

func handler(req events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
	method := strings.ToLower(req.RequestContext.HTTP.Method)

	var res events.APIGatewayProxyResponse
	var err error
	if method == "get" {
		if req.PathParameters["userID"] != "" {
			res, err = getProfileByID(req)
		} else {
			res, err = getProfile(req)
		}

	} else if method == "post" {
		res, err = createProfile(req)
	} else if method == "put" {
		res, err = updateProfile(req)
	} else {
		jsonData, _ := json.Marshal(map[string]interface{}{"message": "Not implemented"})
		res = events.APIGatewayProxyResponse{
			StatusCode: http.StatusNotImplemented,
			Body:       string(jsonData),
		}
	}

	return res, err
}

func main() {
	lambda.Start(handler)
}
