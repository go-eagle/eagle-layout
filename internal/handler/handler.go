package handler

import (
	"github.com/google/wire"

	v1 "github.com/go-eagle/eagle-layout/internal/handler/v1"
)

// here you can add other sets
// ...

// HandlerSet compose all services and all subsets.
var HandlerSet = wire.NewSet(
	v1.NewLoginHandler,
	v1.NewRegisterHandler,
)

type Handler struct {
	Login    *v1.LoginHandler
	Register *v1.RegisterHandler
}

var (
	// Handle expose the handler to outside
	Handle *Handler
)
