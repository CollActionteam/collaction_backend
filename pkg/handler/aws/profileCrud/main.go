package main

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/CollActionteam/collaction_backend/auth"
	"github.com/CollActionteam/collaction_backend/internal/models"
	"github.com/CollActionteam/collaction_backend/internal/profile"
	"github.com/CollActionteam/collaction_backend/pkg/repository"
	"github.com/CollActionteam/collaction_backend/pkg/repository/aws"
	"github.com/CollActionteam/collaction_backend/utils"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type ProfileHandler struct {
	service profile.Service
}

func NewContactHandler() *ProfileHandler {
	profileRepository := repository.NewProfile(aws.NewDynamoConn())
	return &ProfileHandler{service: profile.NewProfileCrudService(profileRepository)}
}

func (h *ProfileHandler) getProfile(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
	userID := req.PathParameters["userID"]
	if userID == "" {
		return utils.GetDataHttpResponse(http.StatusBadRequest, "no profile selected", ""), nil
	}

	profileData, err := h.service.GetProfile(ctx, userID)
	if err != nil {
		return utils.GetDataHttpResponse(http.StatusInternalServerError, "Error Retrieving Profile", ""), nil
	}

	if profileData == nil {
		return utils.GetDataHttpResponse(http.StatusNotFound, "no user Profile found", ""), nil
	}

	return utils.GetDataHttpResponse(http.StatusOK, "Successfully Retrieving Profile", profileData), nil
}

func (h *ProfileHandler) createProfile(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
	usrInf, err := auth.ExtractUserInfo(req)
	if err != nil {
		return utils.GetMessageHttpResponse(http.StatusUnauthorized, err.Error()), nil
	}

	us := models.NewUserInfo(usrInf.UserID(), usrInf.Name(), usrInf.PhoneNumber())
	requestData, err := ValidateProfileRequestData(req, "create")
	if err != nil {
		return utils.GetMessageHttpResponse(http.StatusBadRequest, err.Error()), nil
	}

	err = h.service.CreateProfile(ctx, *us, requestData)
	if err != nil {
		return utils.GetMessageHttpResponse(http.StatusInternalServerError, err.Error()), nil
	}

	return utils.GetMessageHttpResponse(http.StatusOK, "Profile Created"), nil
}

func (h *ProfileHandler) updateProfile(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
	usrInf, err := auth.ExtractUserInfo(req)
	if err != nil {
		return utils.GetMessageHttpResponse(http.StatusUnauthorized, err.Error()), nil
	}

	us := models.NewUserInfo(usrInf.UserID(), usrInf.Name(), usrInf.PhoneNumber())
	requestData, err := ValidateProfileRequestData(req, "update")
	if err != nil {
		return utils.GetMessageHttpResponse(http.StatusBadRequest, err.Error()), nil
	}

	err = h.service.UpdateProfile(ctx, *us, requestData)
	if err != nil {
		return utils.GetMessageHttpResponse(http.StatusInternalServerError, err.Error()), nil
	}

	return utils.GetMessageHttpResponse(http.StatusOK, "profile update successful"), nil
}

func handler(ctx context.Context, req events.APIGatewayV2HTTPRequest) (res events.APIGatewayProxyResponse, err error) {
	method := strings.ToLower(req.RequestContext.HTTP.Method)

	switch method {
	case "get":
		res, err = NewContactHandler().getProfile(ctx, req)
	case "post":
		res, err = NewContactHandler().createProfile(ctx, req)
	case "put":
		res, err = NewContactHandler().updateProfile(ctx, req)
	default:
		res = utils.GetMessageHttpResponse(http.StatusNotImplemented, "Not implemented")
	}
	return
}

func main() {
	lambda.Start(handler)
}

func ValidateProfileRequestData(req events.APIGatewayV2HTTPRequest, method string) (models.Profile, error) {
	var profiledata models.Profile
	err := json.Unmarshal([]byte(req.Body), &profiledata)
	if err != nil {
		return profiledata, err
	}

	err = profiledata.ValidateProfileStruct(method)
	if err != nil {
		return profiledata, err
	}

	return profiledata, nil
}
