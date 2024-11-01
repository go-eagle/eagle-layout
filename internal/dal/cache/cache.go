package cache

import (
	"github.com/go-eagle/eagle/pkg/redis"
	"github.com/google/wire"
)

const (
	// prefix product line prefix
	// you can change it to your custom prefix
	// or set it to empty string if you don't want to use prefix
	prefix = "eagle:"
)

// ProviderSet is cache providers.
var ProviderSet = wire.NewSet(redis.Init, NewUserCache)
