package dal

import (
	"context"

	"github.com/go-eagle/eagle/pkg/log"
	"github.com/go-eagle/eagle/pkg/storage/orm"
	"gorm.io/gorm"

	"github.com/go-eagle/eagle-layout/internal/dal/db/dao"
)

var (
	// DB define a gloabl db
	DB *gorm.DB
)

type DBClient struct {
	db *gorm.DB
}

// Init init db
func Init() (*DBClient, func(), error) {
	// new gorm logger writer with custom filename
	orm.Logger = log.New(log.WithFilename("mysql"))

	// new db
	// 支持多个数据库，配置在 config/{env}/database.yaml 中
	err := orm.New([]string{"default"}...)
	if err != nil {
		return nil, nil, err
	}

	// get default db
	DB, err = orm.GetDB("default")
	if err != nil {
		return nil, nil, err
	}
	sqlDB, err := DB.DB()
	if err != nil {
		return nil, nil, err
	}

	dao.SetDefault(DB)

	cleanFunc := func() {
		sqlDB.Close()
	}
	return &DBClient{db: DB}, cleanFunc, nil
}

func (c *DBClient) GetDB() *gorm.DB {
	return c.db
}

type contextTxKey struct{}

// ExecTx gorm Transaction
func (c *DBClient) ExecTx(ctx context.Context, fn func(ctx context.Context) error) error {
	return c.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		ctx = context.WithValue(ctx, contextTxKey{}, tx)
		return fn(ctx)
	})
}

func (c *DBClient) DBTx(ctx context.Context) *gorm.DB {
	tx, ok := ctx.Value(contextTxKey{}).(*gorm.DB)
	if ok {
		return tx
	}
	return c.db
}
