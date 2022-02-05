package service

import (
	"context"
	"fmt"
)

type IGreeterService interface {
	SayHi(ctx context.Context, name string) (string, error)
}

type greeterService struct {
}

var _ IGreeterService = (*greeterService)(nil)

func newGreeterService(svc *service) *greeterService {
	return &greeterService{}
}

func (s *greeterService) SayHi(ctx context.Context, name string) (string, error) {
	return fmt.Sprintf("Hi %s", name), nil
}
