package repository

import (
	"github.com/go-eagle/eagle-layout/internal/dal"
	"github.com/go-eagle/eagle-layout/internal/dal/cache"
	"github.com/google/wire"
)

// RepositorySet is repo providers.
var RepositorySet = wire.NewSet(
	dal.Init,
	NewUserRepo,
	cache.CacheSet,
)
