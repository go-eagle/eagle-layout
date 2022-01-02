package server

import (
	v1 "github.com/go-eagle/eagle-layout/api/helloworld/greeter/v1"
	"github.com/go-eagle/eagle-layout/internal/routers"
	"github.com/go-eagle/eagle-layout/internal/service"
	"github.com/go-eagle/eagle/pkg/app"
	"github.com/go-eagle/eagle/pkg/transport/http"
)

// NewHTTPServer creates a HTTP server
func NewHTTPServer(c *app.ServerConfig) *http.Server {
	router := routers.NewRouter()

	srv := http.NewServer(
		http.WithAddress(c.Addr),
		http.WithReadTimeout(c.ReadTimeout),
		http.WithWriteTimeout(c.WriteTimeout),
	)

	srv.Handler = router

	v1.RegisterGreeterHTTPServer(router, &service.GreeterGRPCService{})

	return srv
}
