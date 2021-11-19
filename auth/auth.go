package auth

import (
	"errors"

	"github.com/aws/aws-lambda-go/events"
)

type userInfo struct {
	userID      string
	name        string
	phoneNumber string
}

func (usrInf userInfo) UserID() string {
	return usrInf.userID
}

func (usrInf userInfo) Name() string {
	return usrInf.name
}

func (usrInf userInfo) PhoneNumber() string {
	return usrInf.phoneNumber
}

func extractUserInfoFromClaims(claims map[string]string) (usrInf *userInfo) {
	usrInf = &userInfo{}
	usrInf.userID = claims["user_id"]
	usrInf.name = claims["name"]
	usrInf.phoneNumber = claims["phone_number"]
	return
}

func ExtractUserInfo(req events.APIGatewayV2HTTPRequest) (usrInf *userInfo, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("could not extract user info from request")
			usrInf = nil
		}
	}()
	claims := req.RequestContext.Authorizer.JWT.Claims
	usrInf = extractUserInfoFromClaims(claims)
	return
}
