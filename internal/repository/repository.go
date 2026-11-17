package repository

import (
	"github.com/google/wire"

	"github.com/go-eagle/eagle-layout/internal/dal"
	"github.com/go-eagle/eagle-layout/internal/dal/cache"
)

// RepositorySet is repo providers.
var RepositorySet = wire.NewSet(
	dal.Init,
	NewUserRepo,
	cache.CacheSet,
)
