package dal

import (
	"context"

	"github.com/go-eagle/eagle/pkg/log"
	"github.com/go-eagle/eagle/pkg/storage/orm"
	"gorm.io/gorm"

	"github.com/go-eagle/eagle-layout/internal/dal/db/dao"
)

var (
	// DB is the process-wide gorm handle used by generated DAO code.
	DB *gorm.DB
)

// DBClient is the database handle shared by repositories.
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
		if closeErr := sqlDB.Close(); closeErr != nil {
			log.Errorf("close db: %v", closeErr)
		}
	}
	return &DBClient{db: DB}, cleanFunc, nil
}

// GetDB returns the underlying gorm DB.
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

// DBTx returns the transaction stored on ctx, or the default DB.
func (c *DBClient) DBTx(ctx context.Context) *gorm.DB {
	tx, ok := ctx.Value(contextTxKey{}).(*gorm.DB)
	if ok {
		return tx
	}
	return c.db
}
