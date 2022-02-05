package repository

import (
	"github.com/google/wire"
	"gorm.io/gorm"

	"github.com/go-eagle/eagle-layout/internal/model"
)

// ProviderSet is repo providers.
var ProviderSet = wire.NewSet(NewGORMClient)

func NewGORMClient() *gorm.DB {
	return model.GetDB()
}
