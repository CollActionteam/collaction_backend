package crowdaction_test

import (
	"context"
	"fmt"
	"testing"

	cwd "github.com/CollActionteam/collaction_backend/internal/crowdactions"
	"github.com/CollActionteam/collaction_backend/pkg/mocks/repository"
	"github.com/stretchr/testify/assert"
)

func TestCrowdaction_GetCrowdactionById(t *testing.T, ctx context.Context) {
	as := assert.New(t)
	dynamoRepository := &repository.Dynamo{}
	crowdactionID := "Helloworld2"

	t.Run("dev stage", func(t *testing.T) {
		dynamoRepository.On("Send", context.Background(), crowdactionID).Return(nil).Once()

		service := cwd.NewCrowdactionService(dynamoRepository)

		crowdaction, err := service.GetCrowdactionById(ctx, crowdactionID)

		fmt.Printf("Hello world", crowdaction)

		as.NoError(err)

		dynamoRepository.AssertExpectations(t)
	})
}
