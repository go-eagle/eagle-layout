package repository

import (
	"github.com/go-eagle/eagle-layout/internal/dal"
	"github.com/google/wire"
)

// ProviderSet is repo providers.
var ProviderSet = wire.NewSet(dal.Init, NewUserRepo)
