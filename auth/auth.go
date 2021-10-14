package auth

import (
	"github.com/aws/aws-lambda-go/events"
)

type userInfo struct {
	UserID      string
	Name        string
	PhoneNumber string
}

func ExtractUserInfo(request events.APIGatewayProxyRequest) *userInfo {
	var u *userInfo = nil
	if authJWT, hasAuthJWT := request.RequestContext.Authorizer["jwt"]; hasAuthJWT {
		claims := authJWT.(map[string]map[string]interface{})["claims"]
		u = &userInfo{}
		u.UserID = claims["user_id"].(string)
		u.Name = claims["name"].(string)
		u.PhoneNumber = claims["phone_number"].(string)
	}
	return u
}
