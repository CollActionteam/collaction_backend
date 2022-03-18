package auth

import (
	"errors"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

type emailAndVerificationStatus struct {
	email    string
	verified bool
}

type userInfo struct {
	userID      string
	name        *string
	phoneNumber *string
	email       *emailAndVerificationStatus
}

func (usrInf userInfo) UserID() string {
	return usrInf.userID
}

func (usrInf userInfo) Name() *string {
	return usrInf.name
}

func (usrInf userInfo) PhoneNumber() *string {
	return usrInf.phoneNumber
}

func (usrInf userInfo) Email() *emailAndVerificationStatus {
	return usrInf.email
}

func extractUserInfoFromClaims(claims map[string]string) (usrInf *userInfo) {
	usrInf = &userInfo{}
	usrInf.userID = claims["user_id"]
	usrInf.name = nil
	if name, ok := claims["name"]; ok {
		usrInf.name = &name
	}
	usrInf.phoneNumber = nil
	if phoneNumber, ok := claims["phone_number"]; ok {
		usrInf.phoneNumber = &phoneNumber
	}
	usrInf.email = nil
	if email, ok := claims["email"]; ok {
		verified := strings.ToLower(claims["email_verified"]) == "true"
		usrInf.email = &emailAndVerificationStatus{
			email,
			verified,
		}
	}
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
