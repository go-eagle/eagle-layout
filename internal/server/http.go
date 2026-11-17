package server

import (
	"github.com/go-eagle/eagle/pkg/app"
	"github.com/go-eagle/eagle/pkg/transport/http"

	userv1 "github.com/go-eagle/eagle-layout/api/user/v1"
	"github.com/go-eagle/eagle-layout/internal/handler"
	"github.com/go-eagle/eagle-layout/internal/routers"
	"github.com/go-eagle/eagle-layout/internal/service"
)

// NewHTTPServer creates a HTTP server.
// userSvc serves the proto HTTP API. h serves the gin auth routes.
func NewHTTPServer(c *app.Config, userSvc *service.UserServiceServer, h *handler.Handler) *http.Server {
	router := routers.NewRouter(h)

	srv := http.NewServer(
		http.WithAddress(c.HTTP.Addr),
		http.WithReadTimeout(c.HTTP.ReadTimeout),
		http.WithWriteTimeout(c.HTTP.WriteTimeout),
	)

	srv.Handler = router

	userv1.RegisterUserServiceHTTPServer(router, userSvc)

	return srv
}
