package server

import (
	"github.com/google/wire"
)

// ProviderSet is server providers.
// if you want to add rabbitmq, you can append NewRabbitmqConsumerServer in NewSet
var ProviderSet = wire.NewSet(NewHTTPServer, NewGRPCServer, NewRedisConsumerServer)
