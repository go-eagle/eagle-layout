package handler

import (
	"github.com/google/wire"

	v1 "github.com/go-eagle/eagle-layout/internal/handler/v1"
)

// HandlerSet provides HTTP handlers. Wire injects the aggregate into the router.
var HandlerSet = wire.NewSet(
	v1.NewLoginHandler,
	v1.NewRegisterHandler,
	wire.Struct(new(Handler), "*"),
)

// Handler aggregates HTTP handlers constructed by Wire.
type Handler struct {
	Login    *v1.LoginHandler
	Register *v1.RegisterHandler
}
