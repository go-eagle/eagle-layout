package model

import (
	"github.com/go-eagle/eagle/pkg/config"
	"github.com/go-eagle/eagle/pkg/storage/orm"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB
)

// Init init db
func Init() (*gorm.DB, func(), error) {
	cfg, err := loadConf()
	if err != nil {
		return nil, nil, err
	}

	DB = orm.NewMySQL(cfg)
	sqlDB, err := DB.DB()
	if err != nil {
		return nil, nil, err
	}
	cleanFunc := func() {
		sqlDB.Close()
	}
	return DB, cleanFunc, nil
}

// loadConf load gorm config
func loadConf() (ret *orm.Config, err error) {
	var cfg orm.Config
	if err := config.Load("database", &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
