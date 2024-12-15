package server

import (
	"github.com/go-eagle/eagle/pkg/app"
	"github.com/go-eagle/eagle/pkg/transport/http"

	userv1 "github.com/go-eagle/eagle-layout/api/user/v1"
	"github.com/go-eagle/eagle-layout/internal/routers"
	"github.com/go-eagle/eagle-layout/internal/service"
)

// NewHTTPServer creates a HTTP server
// grpc -> if open http by protocol, then add second param: userSvc userv1.UserServiceHTTPServer
func NewHTTPServer(c *app.Config, userSvc *service.UserServiceServer) *http.Server {
	router := routers.NewRouter()

	srv := http.NewServer(
		http.WithAddress(c.HTTP.Addr),
		http.WithReadTimeout(c.HTTP.ReadTimeout),
		http.WithWriteTimeout(c.HTTP.WriteTimeout),
	)

	srv.Handler = router

	userv1.RegisterUserServiceHTTPServer(router, userSvc)

	return srv
}
