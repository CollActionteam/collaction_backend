package repository

import "github.com/stretchr/testify/mock"

type ConfigManager struct {
	mock.Mock
}

func (s *ConfigManager) GetParameter(name string) (string, error) {
	args := s.Mock.Called(name)
	return args.String(0), args.Error(1)
}
