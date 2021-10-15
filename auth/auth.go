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

func ExtractUserInfo(request events.APIGatewayProxyRequest) (usrInf *userInfo, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("could not extract user info from request")
			usrInf = nil
		}
	}()
	if authJWT, hasAuthJWT := request.RequestContext.Authorizer["jwt"]; hasAuthJWT {
		claims := authJWT.(map[string]interface{})["claims"].(map[string]interface{})
		usrInf = &userInfo{}
		usrInf.userID = claims["user_id"].(string)
		usrInf.name = claims["name"].(string)
		usrInf.phoneNumber = claims["phone_number"].(string)
	}
	return
}
