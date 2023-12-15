//go:build wireinject
// +build wireinject

package main

import (
	"github.com/go-eagle/eagle-layout/internal/server"
	"github.com/go-eagle/eagle-layout/internal/tasks"
	eagle "github.com/go-eagle/eagle/pkg/app"
	"github.com/google/wire"
)

func InitApp(cfg *eagle.Config, config *eagle.ServerConfig, tc *tasks.Config) (*eagle.App, func(), error) {
	wire.Build(server.ProviderSet, newApp)
	return &eagle.App{}, nil, nil
}
