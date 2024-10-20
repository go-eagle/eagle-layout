//go:build wireinject
// +build wireinject

package main

import (
	"github.com/go-eagle/eagle-layout/internal/server"

	//"github.com/go-eagle/eagle-layout/internal/service"
	"github.com/go-eagle/eagle-layout/internal/tasks"
	eagle "github.com/go-eagle/eagle/pkg/app"
	logger "github.com/go-eagle/eagle/pkg/log"
	redisMQ "github.com/go-eagle/eagle/pkg/transport/consumer/redis"
	httpSrv "github.com/go-eagle/eagle/pkg/transport/http"
	"github.com/google/wire"
)

func InitApp(cfg *eagle.Config, config *eagle.ServerConfig, tc *tasks.Config) (*eagle.App, func(), error) {
	wire.Build(server.ProviderSet, newApp)
	return &eagle.App{}, nil, nil
}

// 第三个参数需要根据使用的server 进行调整
// 默认使用 redis, 如果使用 rabbitmq 可以改为: rs *rabbitmq.Server
// 然后执行 wire
func newApp(cfg *eagle.Config, hs *httpSrv.Server, rs *redisMQ.Server) *eagle.App {
	return eagle.New(
		eagle.WithName(cfg.Name),
		eagle.WithVersion(cfg.Version),
		eagle.WithLogger(logger.GetLogger()),
		eagle.WithServer(
			// init HTTP server
			hs,
			// init consumer server
			rs,
		),
	)
}
