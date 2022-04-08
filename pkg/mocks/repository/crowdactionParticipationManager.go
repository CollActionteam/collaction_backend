package repository

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type CrowdactionParticipations struct {
	mock.Mock
}

func (e *CrowdactionParticipations) ChangeCrowdactionParticipantCountBy(ctx context.Context, crowdactionID string, count int) error {
	args := e.Called(ctx, crowdactionID, count)
	return args.Error(0)
}
