package service

import (
	"github.com/go-eagle/eagle-layout/internal/repository"
)

// Svc global var
var Svc Service

// Service define all service
type Service interface {
}

// service struct
type service struct {
	repo repository.Repository
}

// New init service
func New(repo repository.Repository) Service {
	return &service{
		repo: repo,
	}
}
