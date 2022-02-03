package repository

import (
	"context"
	"github.com/CollActionteam/collaction_backend/internal/models"
	"github.com/stretchr/testify/mock"
)

type Email struct {
	mock.Mock
}

func (e *Email) Send(ctx context.Context, data models.EmailData) error {
	args := e.Called(ctx, data)
	return args.Error(0)
}
