package crowdaction

import (
	"context"

	m "github.com/CollActionteam/collaction_backend/internal/models"
	"github.com/CollActionteam/collaction_backend/utils"
)

type Service interface {
	GetCrowdactionById(ctx context.Context, crowdactionId string) (*m.CrowdactionData, error)
	GetCrowdactionsByStatus(ctx context.Context, status string, startFrom *utils.PrimaryKey) ([]m.CrowdactionData, error)
	RegisterCrowdaction(ctx context.Context, payload m.CrowdactionData) error
}
type CrowdactionManager interface {
	GetById(pk string, crowdactionId string) (*m.CrowdactionData, error)
	GetByStatus(filterCond string, startFrom *utils.PrimaryKey) ([]m.CrowdactionData, error)
	Register(ctx context.Context, payload m.CrowdactionData) error
}

const (
	KeyDateStart      = "date_start"
	KeyDateEnd        = "date_end"
	KeyDateJoinBefore = "date_limit_join"
)

type crowdactionService struct {
	crowdactionRepository CrowdactionManager
}

func NewCrowdactionService(crowdactionRepository CrowdactionManager) Service {
	return &crowdactionService{crowdactionRepository: crowdactionRepository}
}

func (e *crowdactionService) GetCrowdactionById(ctx context.Context, crowdactionID string) (*m.CrowdactionData, error) {
	return e.crowdactionRepository.GetById(utils.PKCrowdaction, crowdactionID)
}

func (e *crowdactionService) GetCrowdactionsByStatus(ctx context.Context, status string, startFrom *utils.PrimaryKey) ([]m.CrowdactionData, error) {
	return e.crowdactionRepository.GetByStatus(status, startFrom)
}

func (e *crowdactionService) RegisterCrowdaction(ctx context.Context, payload m.CrowdactionData) error {
	if err := e.crowdactionRepository.Register(ctx, payload); err != nil {
		return err
	}
	return nil // should this be ultimately be nil?
}
