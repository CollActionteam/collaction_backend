package participation

import (
	"context"

	"github.com/CollActionteam/collaction_backend/internal/constants"
	m "github.com/CollActionteam/collaction_backend/internal/models"
	"github.com/CollActionteam/collaction_backend/models"
	"github.com/CollActionteam/collaction_backend/utils"
	"github.com/pkg/errors"
)

type Service interface {
	GetParticipation(ctx context.Context, userID string, crowdactionID string) (*m.ParticipationRecord, error)
	RegisterParticipation(ctx context.Context, userID string, name string, crowdaction *models.Crowdaction, payload m.JoinPayload) error
	CancelParticipation(ctx context.Context, userID string, crowdaction *models.Crowdaction) error
}

type ParticipationManager interface {
	Get(ctx context.Context, userID string, crowdactionID string) (*m.ParticipationRecord, error)
	Register(ctx context.Context, userID string, name string, crowdaction *models.Crowdaction, payload m.JoinPayload) error
	Cancel(ctx context.Context, userID string, crowdaction *models.Crowdaction) error
}

type participationService struct {
	participationRepository ParticipationManager
}

// TODO move to handler
func NewParticipationService(participationRepository ParticipationManager) Service {
	return &participationService{
		participationRepository: participationRepository,
	}
}

func recordEvent(userID string, crowdactionID string, commitments []string, count int) error {
	qc := utils.CreateQueueClient()
	event := m.ParticipationEvent{
		UserID:        userID,
		CrowdactionID: crowdactionID,
		Commitments:   commitments,
		Count:         count,
	}
	return utils.SendQueueMessage(qc, constants.ParticipationQueueName, event)
}

func (s *participationService) GetParticipation(ctx context.Context, userID string, crowdactionID string) (*m.ParticipationRecord, error) {
	return s.participationRepository.Get(ctx, userID, crowdactionID)
}

func (s *participationService) RegisterParticipation(ctx context.Context, userID string, name string, crowdaction *models.Crowdaction, payload m.JoinPayload) error {
	if !utils.IsFutureDateString(crowdaction.DateLimitJoin) {
		return errors.New("cannot change participation for this crowdaction anymore")
	}
	if err := s.participationRepository.Register(ctx, userID, name, crowdaction, payload); err != nil {
		return err
	}
	return recordEvent(userID, crowdaction.CrowdactionID, payload.Commitments, +1)
}

func (s *participationService) CancelParticipation(ctx context.Context, userID string, crowdaction *models.Crowdaction) error {
	if !utils.IsFutureDateString(crowdaction.DateEnd) {
		return errors.New("cannot change participation for this crowdaction anymore")
	}

	part, err := s.participationRepository.Get(ctx, userID, crowdaction.CrowdactionID)
	if err != nil {
		return err
	}
	if part == nil {
		return errors.New("not participating")
	}
	if err = s.participationRepository.Cancel(ctx, userID, crowdaction); err != nil {
		return err
	}
	return recordEvent(userID, crowdaction.CrowdactionID, part.Commitments, -1)
}
