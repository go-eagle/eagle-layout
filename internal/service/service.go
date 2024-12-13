package service

import (
	"github.com/go-eagle/eagle-layout/internal/repository"
	"github.com/google/wire"
)

// ServiceSet is service providers.
var ServiceSet = wire.NewSet(
	NewUserServiceServer,
	repository.RepositorySet,
)
