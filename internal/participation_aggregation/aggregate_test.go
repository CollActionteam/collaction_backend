package participation_aggregation_test

import (
	"context"
	"testing"

	"github.com/CollActionteam/collaction_backend/internal/participation_aggregation"
	"github.com/CollActionteam/collaction_backend/pkg/mocks/repository"
	"github.com/stretchr/testify/assert"
)

func TestContact_ChangeCrowdactionParticipantCountBy(t *testing.T) {
	as := assert.New(t)
	participationManager := &repository.CrowdactionParticipations{}
	count := 100
	crowdactionID := "Some crowdaction"

	t.Run("increase participation", func(t *testing.T) {
		service := participation_aggregation.NewParticipationAggregationService(participationManager)

		participationManager.On("ChangeCrowdactionParticipantCountBy", context.Background(), crowdactionID, count).Return(nil).Once()

		err := service.ChangeCrowdactionParticipantCountBy(context.Background(), crowdactionID, count)
		as.NoError(err)

		participationManager.AssertExpectations(t)
	})

	t.Run("decrease participation", func(t *testing.T) {
		service := participation_aggregation.NewParticipationAggregationService(participationManager)

		participationManager.On("ChangeCrowdactionParticipantCountBy", context.Background(), crowdactionID, count*-1).Return(nil).Once()

		err := service.ChangeCrowdactionParticipantCountBy(context.Background(), crowdactionID, count*-1)
		as.NoError(err)

		participationManager.AssertExpectations(t)
	})

}
