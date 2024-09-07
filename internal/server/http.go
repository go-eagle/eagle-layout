package server

import (
	v1 "github.com/go-eagle/eagle-layout/api/helloworld/greeter/v1"
	"github.com/go-eagle/eagle-layout/internal/routers"
	"github.com/go-eagle/eagle-layout/internal/service"
	"github.com/go-eagle/eagle/pkg/app"
	"github.com/go-eagle/eagle/pkg/transport/http"
)

// NewHTTPServer creates a HTTP server
func NewHTTPServer(c *app.Config, greeterSvc *service.GreeterServiceServer) *http.Server {
	router := routers.NewRouter()

	srv := http.NewServer(
		http.WithAddress(c.HTTP.Addr),
		http.WithReadTimeout(c.HTTP.ReadTimeout),
		http.WithWriteTimeout(c.HTTP.WriteTimeout),
	)

	srv.Handler = router

	v1.RegisterGreeterServiceHTTPServer(router, greeterSvc)

	return srv
}
