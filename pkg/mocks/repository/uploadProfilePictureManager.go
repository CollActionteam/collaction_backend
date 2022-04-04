package repository

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type ProfilePicture struct {
	mock.Mock
}

func (p *ProfilePicture) GetUploadUrl(ctx context.Context, ext string, userID string) (string, error) {
	outputs := p.Mock.Called(ctx, ext, userID)

	uploadUrl := outputs.String(0)
	err := outputs.Error(1)

	return uploadUrl, err
}
