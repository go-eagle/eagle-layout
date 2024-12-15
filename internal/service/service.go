package service

import (
	"github.com/google/wire"

	"github.com/go-eagle/eagle-layout/internal/repository"
)

// ServiceSet is service providers.
var ServiceSet = wire.NewSet(
	NewUserServiceServer, // for grpc(inlucde http from grpc)
	NewUserService,       // for http
	repository.RepositorySet,
)
