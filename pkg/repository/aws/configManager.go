package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/pkg/errors"
)

type ConfigManager struct {
	Client *ssm.SSM
}

func NewConfigManager(session *session.Session) *ConfigManager {
	return &ConfigManager{Client: ssm.New(session)}
}

func (s *ConfigManager) GetParameter(name string) (string, error) {
	param, err := s.Client.GetParameter(&ssm.GetParameterInput{Name: aws.String(name)})
	if err != nil {
		return "", errors.WithStack(err)
	}

	return *param.Parameter.Value, nil
}
