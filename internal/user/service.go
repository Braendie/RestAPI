package user

import (
	"context"

	"github.com/Braendie/RestAPI/pkg/logging"
)

type Service struct {
	storage Storage
	logger *logging.Logger
}

func (s *Service) Create(ctx context.Context, dto CreateUserDTO) (u User, err error) {
	//Todo to next one
	return
}