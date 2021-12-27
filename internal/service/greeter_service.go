package service

import (
	"context"
	"fmt"

	"github.com/go-eagle/eagle-layout/internal/repository"
)

type GreeterService interface {
	SayHi(ctx context.Context, name string) (string, error)
}

type greeterService struct {
	repo repository.Repository
}

var _ GreeterService = (*greeterService)(nil)

func newGreeterService(svc *service) *greeterService {
	return &greeterService{repo: svc.repo}
}

func (s *greeterService) SayHi(ctx context.Context, name string) (string, error) {
	return fmt.Sprintf("Hi %s", name), nil
}
