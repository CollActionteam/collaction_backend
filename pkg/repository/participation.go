package repository

import (
	"context"
	"errors"

	"github.com/CollActionteam/collaction_backend/internal/constants"
	m "github.com/CollActionteam/collaction_backend/internal/models"
	"github.com/CollActionteam/collaction_backend/models"
	"github.com/CollActionteam/collaction_backend/pkg/repository/aws"
	"github.com/CollActionteam/collaction_backend/utils"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Participation interface {
	Get(ctx context.Context, userID string, crowdactionID string) (*m.ParticipationRecord, error)
	Register(ctx context.Context, userID string, name string, crowdaction *models.Crowdaction, payload m.JoinPayload) error
	Cancel(ctx context.Context, userID string, crowdaction *models.Crowdaction) error
	List(ctx context.Context, userID string) (*[]m.ParticipationRecord, error)
}

type participation struct {
	dbClient *aws.Dynamo
}

func NewParticipation(dynamo *aws.Dynamo) Participation {
	return &participation{
		dbClient: dynamo,
	}
}
func (s *participation) Get(ctx context.Context, userID string, crowdactionID string) (*m.ParticipationRecord, error) {
	pk := utils.PrefixParticipationPK_UserID + userID
	sk := utils.PrefixParticipationSK_CrowdactionID + crowdactionID
	item, err := s.dbClient.GetDBItem(constants.TableName, pk, sk)
	if item == nil || err != nil {
		return nil, err
	}
	var r m.ParticipationRecord
	err = dynamodbattribute.UnmarshalMap(item, &r)
	return &r, err
}

func (s *participation) Register(ctx context.Context, userID string, name string, crowdaction *models.Crowdaction, payload m.JoinPayload) error {
	/* TODO Password not required when joining for MVP
	if crowdaction.PasswordJoin != "" && crowdaction.PasswordJoin != payload.Password {
		return fmt.Errorf("invalid password")
	}
	*/
	part, err := s.Get(ctx, userID, crowdaction.CrowdactionID)
	if part != nil {
		err = errors.New("already participating")
	}
	if err != nil {
		return err
	}
	pk := utils.PrefixParticipationPK_UserID + userID
	sk := utils.PrefixParticipationSK_CrowdactionID + crowdaction.CrowdactionID
	return s.dbClient.PutDBItem(constants.TableName, pk, sk, m.ParticipationRecord{
		UserID:        userID,
		Name:          name,
		CrowdactionID: crowdaction.CrowdactionID,
		Title:         crowdaction.Title,
		Commitments:   payload.Commitments,
		Date:          utils.GetDateStringNow(),
	})
}

func (s *participation) Cancel(ctx context.Context, userID string, crowdaction *models.Crowdaction) error {
	pk := utils.PrefixParticipationPK_UserID + userID
	sk := utils.PrefixParticipationSK_CrowdactionID + crowdaction.CrowdactionID
	return s.dbClient.DeleteDBItem(constants.TableName, pk, sk)
}

func (s *participation) List(ctx context.Context, userID string) (*[]m.ParticipationRecord, error) {
	pk := utils.PrefixParticipationPK_UserID + userID
	//s.dbClient.
}
