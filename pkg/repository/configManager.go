package repository

import "github.com/pkg/errors"

var myConfig = map[string]string{
	"/collaction/dev/contact/email":        "dev@gmail.com",
	"/collaction/debug/contact/email":      "dev@gmail.com",
	"/collaction/production/contact/email": "production@gmail.com",
}

type ConfigManager struct{}

func NewConfigManager() *ConfigManager {
	return &ConfigManager{}
}

func (s *ConfigManager) GetParameter(name string) (string, error) {
	value, ok := myConfig[name]
	if !ok {
		return "", errors.New("key is not found")
	}

	return value, nil
}
