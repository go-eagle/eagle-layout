package handler

import (
	"github.com/google/wire"

	hanlderv1 "github.com/go-eagle/eagle-layout/internal/handler/v1"
)

// ProviderSet is repo providers.
var ProviderSet = wire.NewSet(
	hanlderv1.NewLoginHandler,
)

type Handler struct {
	LoginHandler *hanlderv1.LoginHandler
	// more handlers
}
