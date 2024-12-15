//go:build wireinject
// +build wireinject

package handler

import (
	"github.com/go-eagle/eagle-layout/internal/service"
	"github.com/google/wire"
)

func NewHandler() (*Handler, func(), error) {
	wire.Build(
		service.ServiceSet,
		HandlerSet,
		wire.Struct(new(Handler), "*"),
	)
	return &Handler{}, nil, nil
}
