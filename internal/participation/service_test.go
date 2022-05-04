package participation_test

import (
	"context"
	"testing"

	m "github.com/CollActionteam/collaction_backend/internal/models"
	"github.com/CollActionteam/collaction_backend/internal/participation"
	"github.com/CollActionteam/collaction_backend/models"
	"github.com/CollActionteam/collaction_backend/pkg/mocks/repository"
	"github.com/stretchr/testify/assert"
)

func TestParticipation_RegisterParticipation(t *testing.T) {
	repository := &repository.Participation{}
	oldDate := "1999-02-08"

	t.Run("Given the DateLimitJoin is before today when register participation then should return an error", func(t *testing.T) {
		service := participation.NewParticipationService(repository)
		err := service.RegisterParticipation(context.Background(), "123", "make cool stuffs",
			&models.Crowdaction{DateLimitJoin: oldDate}, m.JoinPayload{})
		assert.EqualError(t, err, "cannot change participation for this crowdaction anymore")
	})
}

func Test_GetParticipationsUser(t *testing.T) {
	ast := assert.New(t)
	repository := &repository.Participation{}
	cc := []m.ParticipationRecord{}

	t.Run("testing GetParticipationsUser", func(t *testing.T) {
		userID := "123"

		repository.On("ListByUser", context.Background(), userID).Return(&cc, nil).Once()

		service := participation.NewParticipationService(repository)
		crowdactions, err := service.GetParticipationsUser(context.Background(), userID)
		ast.NoError(err)
		ast.Equal(&cc, crowdactions)

		repository.AssertExpectations(t)
	})
}

func Test_GetParticipationsCrowdaction(t *testing.T) {
	ast := assert.New(t)
	cc := []m.ParticipationRecord{}
	repository := &repository.Participation{}

	t.Run("testing GetParticipationsCrowdaction", func(t *testing.T) {
		crowdactionID := "1"

		repository.On("ListByCrowdaction", context.Background(), crowdactionID).Return(&cc, nil).Once()

		service := participation.NewParticipationService(repository)
		crowdactions, err := service.GetParticipationsCrowdaction(context.Background(), crowdactionID)
		ast.NoError(err)
		ast.Equal(&cc, crowdactions)

		repository.AssertExpectations(t)
	})
}
