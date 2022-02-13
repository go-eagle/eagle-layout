package cache

import (
	"github.com/go-eagle/eagle/pkg/redis"
	"github.com/google/wire"
)

// ProviderSet is cache providers.
var ProviderSet = wire.NewSet(redis.Init)
