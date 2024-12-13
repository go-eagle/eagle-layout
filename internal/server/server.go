package server

import (
	"github.com/google/wire"
)

// ServerSet is server providers.
// if you want to add rabbitmq, you can append NewRabbitmqConsumerServer in NewSet
var ServerSet = wire.NewSet(
	NewHTTPServer,
	NewGRPCServer,
	NewRedisConsumerServer,
)
