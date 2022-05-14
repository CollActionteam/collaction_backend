package repository

import (
	"context"

	m "github.com/CollActionteam/collaction_backend/internal/models"
	"github.com/CollActionteam/collaction_backend/models"
	"github.com/stretchr/testify/mock"
)

type Participation struct {
	mock.Mock
}

func (s *Participation) Get(ctx context.Context, userID string, crowdactionID string) (*m.ParticipationRecord, error) {
	args := s.Mock.Called(ctx, userID, crowdactionID)
	return args.Get(0).(*m.ParticipationRecord), args.Error(1)
}

func (s *Participation) Register(ctx context.Context, userID string, name string, crowdaction *models.Crowdaction, payload m.JoinPayload) error {
	args := s.Mock.Called(ctx, userID, name, crowdaction, payload)
	return args.Error(0)
}

func (s *Participation) Cancel(ctx context.Context, userID string, crowdaction *models.Crowdaction) error {
	args := s.Mock.Called(ctx, userID, crowdaction)
	return args.Error(0)
}

func (s *Participation) ListByUser(ctx context.Context, userID string) (*[]m.ParticipationRecord, error) {
	args := s.Mock.Called(ctx, userID)
	return args.Get(0).(*[]m.ParticipationRecord), args.Error(1)
}

func (s *Participation) ListByCrowdaction(ctx context.Context, crowdactionID string) (*[]m.ParticipationRecord, error) {
	args := s.Mock.Called(ctx, crowdactionID)
	return args.Get(0).(*[]m.ParticipationRecord), args.Error(1)
}
