package server

import (
	"github.com/google/wire"

	"github.com/go-eagle/eagle-layout/internal/service"
)

// ServerSet is server providers.
// if you want to add rabbitmq, you can append NewRabbitmqConsumerServer in NewSet
var ServerSet = wire.NewSet(
	NewHTTPServer,
	NewGRPCServer,
	NewRedisConsumerServer,
	service.ServiceSet,
)
