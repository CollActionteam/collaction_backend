package contact_test

import (
	"context"
	"fmt"
	"github.com/CollActionteam/collaction_backend/internal/constants"
	"github.com/CollActionteam/collaction_backend/internal/contact"
	"github.com/CollActionteam/collaction_backend/internal/models"
	"github.com/CollActionteam/collaction_backend/pkg/mocks/repository"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestContact_SendEmail(t *testing.T) {
	as := assert.New(t)
	emailRepository := &repository.Email{}
	configManager := &repository.ConfigManager{}

	emailRequest := models.EmailContactRequest{Email: "test@email.com", Subject: "test subject", Message: "test message", AppVersion: "version 1"}

	t.Run("dev stage", func(t *testing.T) {
		stage := "dev"
		recipientKey := fmt.Sprintf(constants.RecipientEmail, stage)
		recipientValue := "dev@email.com"
		service := contact.NewContactService(emailRepository, configManager, stage)

		emailData := models.EmailData{
			Recipient:  recipientValue,
			Message:    fmt.Sprintf(contact.EmailMessageFormat, emailRequest.Message, contact.Separator, emailRequest.AppVersion),
			Subject:    emailRequest.Subject,
			Sender:     emailRequest.Email,
			ReplyEmail: emailRequest.Email,
		}

		emailRepository.On("Send", context.Background(), emailData).Return(nil).Once()
		configManager.On("GetParameter", recipientKey).Return(recipientValue, nil).Once()

		err := service.SendEmail(context.Background(), emailRequest)
		as.NoError(err)

		emailRepository.AssertExpectations(t)
		configManager.AssertExpectations(t)
	})

	t.Run("production stage", func(t *testing.T) {
		stage := "production"
		recipientKey := fmt.Sprintf(constants.RecipientEmail, stage)
		recipientValue := "production@email.com"
		service := contact.NewContactService(emailRepository, configManager, stage)

		emailData := models.EmailData{
			Recipient:  recipientValue,
			Message:    fmt.Sprintf(contact.EmailMessageFormat, emailRequest.Message, contact.Separator, emailRequest.AppVersion),
			Subject:    emailRequest.Subject,
			Sender:     emailRequest.Email,
			ReplyEmail: emailRequest.Email,
		}

		emailRepository.On("Send", context.Background(), emailData).Return(nil).Once()
		configManager.On("GetParameter", recipientKey).Return(recipientValue, nil).Once()

		err := service.SendEmail(context.Background(), emailRequest)
		as.NoError(err)

		emailRepository.AssertExpectations(t)
		configManager.AssertExpectations(t)
	})

}
