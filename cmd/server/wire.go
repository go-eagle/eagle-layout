//go:build wireinject
// +build wireinject

package main

import (
	"github.com/go-eagle/eagle-layout/internal/server"
	eagle "github.com/go-eagle/eagle/pkg/app"
	logger "github.com/go-eagle/eagle/pkg/log"
	"github.com/go-eagle/eagle/pkg/transport/grpc"
	httpSrv "github.com/go-eagle/eagle/pkg/transport/http"
	"github.com/google/wire"
)

func InitApp(cfg *eagle.Config) (*eagle.App, func(), error) {
	// wire.Build(server.ProviderSet, service.ProviderSet, repository.ProviderSet, cache.ProviderSet, newApp)
	wire.Build(server.ProviderSet, newApp)
	return &eagle.App{}, nil, nil
}

func newApp(cfg *eagle.Config, hs *httpSrv.Server, gs *grpc.Server) *eagle.App {
	logger.Init(logger.WithFilename("app"))

	return eagle.New(
		eagle.WithName(cfg.Name),
		eagle.WithVersion(cfg.Version),
		eagle.WithLogger(logger.GetLogger()),
		eagle.WithServer(
			// init HTTP server
			hs,
			// init GRPC server
			gs,
		),
	)
}
