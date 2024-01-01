package service

import (
	"context"

	"github.com/go-eagle/eagle-layout/internal/repository"
)

// GreeterService define a interface
type GreeterService interface {
	Hello(ctx context.Context) error
}

type greeterService struct {
	repo repository.UserRepo
}

var _ GreeterService = (*greeterService)(nil)

func newGreeterService(repo repository.UserRepo) *greeterService {
	return &greeterService{
		repo: repo,
	}
}

// Hello .
func (s *greeterService) Hello(ctx context.Context) error {
	return nil
}
