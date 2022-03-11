package participation_aggregation

import "context"

type ParticipationManager interface {
	ChangeCrowdactionParticipantCountBy(ctx context.Context, crowdactionID string, count int) error
}

type Service interface {
	ChangeCrowdactionParticipantCountBy(ctx context.Context, crowdactionID string, count int) error
}

type aggregate struct {
	participationManager ParticipationManager
}

func NewParticipationAggregationService(participationManager ParticipationManager) Service {
	return &aggregate{participationManager: participationManager}
}

func (e *aggregate) ChangeCrowdactionParticipantCountBy(ctx context.Context, crowdactionID string, count int) error {
	return e.participationManager.ChangeCrowdactionParticipantCountBy(ctx, crowdactionID, count)
}
