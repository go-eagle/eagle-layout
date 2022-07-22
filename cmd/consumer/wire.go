// +build wireinject

package main

import (
	"github.com/go-eagle/eagle-layout/internal/server"
	eagle "github.com/go-eagle/eagle/pkg/app"
	"github.com/google/wire"
)

func InitApp(cfg *eagle.Config, config *eagle.ServerConfig) (*eagle.App, error) {
	wire.Build(server.ProviderSet, newApp)
	return &eagle.App{}, nil
}
