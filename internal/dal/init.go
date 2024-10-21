package dal

import (
	"context"

	"github.com/go-eagle/eagle/pkg/storage/orm"
	"gorm.io/gorm"
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
	// new gorm logger writer with custom filename for sql
	// orm.LogWriter = orm.NewLogWriter(log.New(log.WithFilename("mysql")))

	// new db
	err := orm.New([]string{"default"}...)
	if err != nil {
		return nil, nil, err
	}

	// get first db
	DB, err = orm.GetDB("default")
	if err != nil {
		return nil, nil, err
	}
	sqlDB, err := DB.DB()
	if err != nil {
		return nil, nil, err
	}

	// query.SetDefault(DB)

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
