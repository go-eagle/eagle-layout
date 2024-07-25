package service

import (
	"context"

	"github.com/go-eagle/eagle-layout/internal/repository"
)

// GreeterService define a interface
type GreeterService interface {
	Hello(ctx context.Context, name string) (string, error)
}

// greeterService define a struct
type greeterService struct {
	repo repository.UserRepo
}

var _ GreeterService = (*greeterService)(nil)

// NewGreeterService create a service
func NewGreeterService(repo repository.UserRepo) GreeterService {
	return &greeterService{
		repo: repo,
	}
}

// Hello .
func (s *greeterService) Hello(ctx context.Context, name string) (string, error) {
	return "hello " + name, nil
}
