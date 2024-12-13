package service

import (
	"context"

	"github.com/go-eagle/eagle-layout/internal/repository"
)

// UserService define a interface
type UserService interface {
	Login(ctx context.Context, username, password string) (string, error)
}

// greeterService define a struct
type userService struct {
	repo repository.UserRepo
}

var _ UserService = (*userService)(nil)

// NewUserService create a service
func NewUserService(repo repository.UserRepo) UserService {
	return &userService{
		repo: repo,
	}
}

// Hello .
func (s *userService) Login(ctx context.Context, username, password string) (string, error) {
	return "hello " + username, nil
}
