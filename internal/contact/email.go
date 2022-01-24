package contact

import (
	"context"
	"fmt"
	"github.com/CollActionteam/collaction_backend/internal/constants"
	"github.com/CollActionteam/collaction_backend/internal/models"
)

const separator = "### app version:"

type EmailRepository interface {
	Send(ctx context.Context, data models.EmailData) error
}

type ConfigManager interface {
	GetParameter(name string) (string, error)
}

type Service interface {
	SendEmail(ctx context.Context, data models.EmailContactRequest) error
}

type contact struct {
	emailRepository EmailRepository
	configManager   ConfigManager
	stage           string
}

func NewContactService(emailRepository EmailRepository, configManager ConfigManager, stage string) Service {
	return &contact{emailRepository: emailRepository, configManager: configManager, stage: stage}
}

func (e *contact) SendEmail(ctx context.Context, data models.EmailContactRequest) error {
	recipient, err := e.configManager.GetParameter(fmt.Sprintf(constants.RecipientEmail, e.stage))
	if err != nil {
		return err
	}

	return e.emailRepository.Send(ctx, models.EmailData{
		Recipient:  recipient,
		Message:    fmt.Sprintf("%s %s %s", data.Message, separator, data.AppVersion),
		Subject:    data.Subject,
		Sender:     data.Email,
		ReplyEmail: data.Email,
	})
}
