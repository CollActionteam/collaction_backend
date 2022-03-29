package crowdaction_test

import (
	"context"
	"fmt"
	"testing"

	cwd "github.com/CollActionteam/collaction_backend/internal/crowdactions"
	m "github.com/CollActionteam/collaction_backend/internal/models"
	"github.com/CollActionteam/collaction_backend/pkg/mocks/repository"

	"github.com/CollActionteam/collaction_backend/utils"
	"github.com/stretchr/testify/assert"
)

func TestCrowdaction_GetCrowdactionById(t *testing.T) {
	as := assert.New(t)
	dynamoRepository := &repository.Dynamo{}
	var ctx context.Context
	var crowdactions *m.CrowdactionData
	crowdactionID := "sustainability#food#185f66fd"

	t.Run("dev stage", func(t *testing.T) {
		dynamoRepository.On("GetById", utils.PKCrowdaction, crowdactionID).Return(crowdactions, nil).Once()

		service := cwd.NewCrowdactionService(dynamoRepository)

		crowdaction, err := service.GetCrowdactionById(ctx, crowdactionID)

		fmt.Println("Hello world", crowdaction)

		as.NoError(err)

		dynamoRepository.AssertExpectations(t)
	})
}
