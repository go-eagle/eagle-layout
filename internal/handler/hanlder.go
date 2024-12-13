package handler

import (
	"github.com/google/wire"

	v1 "github.com/go-eagle/eagle-layout/internal/handler/v1"
	"github.com/go-eagle/eagle-layout/internal/service"
)

var authSet = wire.NewSet(
	v1.NewLoginHandler,
)

// here you can add other sets
// ...

// HandlerSet compose all services and all subsets.
var HandlerSet = wire.NewSet(
	service.ServiceSet,
	authSet,
)
