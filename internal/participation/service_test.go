package participation_test

import (
	"context"
	m "github.com/CollActionteam/collaction_backend/internal/models"
	"github.com/CollActionteam/collaction_backend/internal/participation"
	"github.com/CollActionteam/collaction_backend/models"
	"github.com/CollActionteam/collaction_backend/pkg/mocks/repository"
	"github.com/stretchr/testify/assert"
	"testing"
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
