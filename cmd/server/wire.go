//go:build wireinject
// +build wireinject

package main

import (
	eagle "github.com/go-eagle/eagle/pkg/app"
	logger "github.com/go-eagle/eagle/pkg/log"
	"github.com/go-eagle/eagle/pkg/transport/grpc"
	httpSrv "github.com/go-eagle/eagle/pkg/transport/http"
	"github.com/google/wire"

	"github.com/go-eagle/eagle-layout/internal/handler"
	v1 "github.com/go-eagle/eagle-layout/internal/handler/v1"
	"github.com/go-eagle/eagle-layout/internal/server"
)

type Application struct {
	*eagle.App
	LoginHandler *v1.LoginHandler
}

func InitApp(cfg *eagle.Config) (*Application, func(), error) {
	wire.Build(
		server.ServerSet,
		handler.HandlerSet, // 汇总所有 Handler 的依赖
		newApp,
		wire.Struct(new(Application), "*"),
	)
	return &Application{}, nil, nil
}

func newApp(cfg *eagle.Config, hs *httpSrv.Server, gs *grpc.Server) *eagle.App {
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
