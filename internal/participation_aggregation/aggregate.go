package participation_aggregation

import "context"

type CrowdactionParticipationManager interface {
	ChangeCrowdactionParticipantCountBy(ctx context.Context, crowdactionID string, count int) error
}

type Service interface {
	ChangeCrowdactionParticipantCountBy(ctx context.Context, crowdactionID string, count int) error
}

type aggregate struct {
	crowdactionParticipationManager CrowdactionParticipationManager
}

func NewParticipationAggregationService(crowdactionParticipationManager CrowdactionParticipationManager) Service {
	return &aggregate{crowdactionParticipationManager: crowdactionParticipationManager}
}

func (e *aggregate) ChangeCrowdactionParticipantCountBy(ctx context.Context, crowdactionID string, count int) error {
	return e.crowdactionParticipationManager.ChangeCrowdactionParticipantCountBy(ctx, crowdactionID, count)
}
