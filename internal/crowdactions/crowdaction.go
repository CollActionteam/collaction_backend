package crowdaction

import (
	"context"
	"fmt"

	m "github.com/CollActionteam/collaction_backend/internal/models"
	"github.com/CollActionteam/collaction_backend/utils"
)

type Service interface {
	GetCrowdactionById(ctx context.Context, crowdactionId string) (*m.CrowdactionData, error)
	GetCrowdactionsByStatus(ctx context.Context, status string, startFrom *utils.PrimaryKey) ([]m.CrowdactionData, error)
}
type CrowdactionManager interface {
	GetById(pk string, crowdactionId string) (*m.CrowdactionData, error)
	GetByStatus(filterCond string, startFrom *utils.PrimaryKey) ([]m.CrowdactionData, error)
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

/**
	GET Crowdaction by Id
**/
func (e *crowdactionService) GetCrowdactionById(ctx context.Context, crowdactionID string) (*m.CrowdactionData, error) {
	fmt.Println("GetCrowdactionById", crowdactionID)
	return e.crowdactionRepository.GetById(utils.PKCrowdaction, crowdactionID)
}

/**
	GET Crowdaction by Status
**/
func (e *crowdactionService) GetCrowdactionsByStatus(ctx context.Context, status string, startFrom *utils.PrimaryKey) ([]m.CrowdactionData, error) {
	fmt.Println("GetCrowdactionsByStatus", status, startFrom)
	return e.crowdactionRepository.GetByStatus(status, startFrom)
}
