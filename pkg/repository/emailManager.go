package repository

import (
	"context"
	"fmt"
	"github.com/CollActionteam/collaction_backend/internal/models"
)

type Email struct {
	Username string
	Password string
}

func NewEmail(username, password string) *Email {
	return &Email{Username: username, Password: password}
}

func (e *Email) Send(ctx context.Context, data models.EmailData) error {
	fmt.Println("email is sent: ", data)
	return nil
}
