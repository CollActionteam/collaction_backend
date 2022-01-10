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

type Contact struct {
	EmailRepository EmailRepository
	ConfigManager   ConfigManager
	Stage           string
}

func NewContactService(emailRepository EmailRepository, configManager ConfigManager, stage string) *Contact {
	return &Contact{EmailRepository: emailRepository, ConfigManager: configManager, Stage: stage}
}

func (e *Contact) SendEmail(ctx context.Context, data models.EmailContactRequest) error {
	recipient, err := e.ConfigManager.GetParameter(fmt.Sprintf(constants.RecipientEmail, e.Stage))
	if err != nil {
		return err
	}

	return e.EmailRepository.Send(ctx, models.EmailData{
		Recipient:  recipient,
		Message:    fmt.Sprintf("%s %s %s", data.Message, separator, data.AppVersion),
		Subject:    data.Subject,
		Sender:     data.Email,
		ReplyEmail: data.Email,
	})
}
