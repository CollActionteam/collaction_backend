package utils

import (
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

// Deprecated: Use GetDataHttpResponse instead (uses APIGatewayV2HTTPResponse)
func GetMessageHttpResponse(statusCode int, msg string) events.APIGatewayProxyResponse {
	// "Cannot go wrong"
	jsonPayload, _ := json.Marshal(map[string]interface{}{"message": msg})
	return events.APIGatewayProxyResponse{
		Body:       string(jsonPayload),
		StatusCode: statusCode,
	}
}

func GetDataHttpResponse(statusCode int, msg string, data interface{}) events.APIGatewayV2HTTPResponse {
	resp := struct {
		Message string
		Data    interface{}
		Status  int
	}{
		Message: msg,
		Data:    data,
		Status:  statusCode,
	}
	jsonData, _ := json.Marshal(resp)
	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusInternalServerError,
		Body:       string(jsonData),
	}
}

func CreateMessageHttpResponse(statusCode int, msg string) events.APIGatewayV2HTTPResponse {
	jsonPayload, _ := json.Marshal(map[string]interface{}{"message": msg})
	return events.APIGatewayV2HTTPResponse{
		Body:       string(jsonPayload),
		StatusCode: statusCode,
	}
}
