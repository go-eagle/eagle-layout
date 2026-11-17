package service

import (
	"github.com/google/wire"

	"github.com/go-eagle/eagle-layout/internal/repository"
)

// ServiceSet is service providers.
var ServiceSet = wire.NewSet(
	NewUserServiceServer, // gRPC adapter, including the transcoded HTTP API
	NewUserService,       // business service shared by handlers and gRPC
	repository.RepositorySet,
)
