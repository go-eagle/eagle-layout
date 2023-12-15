package repository

import (
	"github.com/go-eagle/eagle-layout/internal/model"
	"github.com/go-eagle/eagle/pkg/redis"
	"github.com/google/wire"
)

// ProviderSet is repo providers.
var ProviderSet = wire.NewSet(model.Init, redis.Init)
