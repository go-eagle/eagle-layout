package service

import (
	"context"
	"fmt"
)

type HiService struct {
}

func NewHiService() *HiService {
	return &HiService{}
}

func (s *HiService) SayHi(ctx context.Context, name string) (string, error) {
	return fmt.Sprintf("Hi %s", name), nil
}
