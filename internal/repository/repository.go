package repository

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"

	"github.com/go-eagle/eagle/pkg/storage/sql"
)

var (
	// ErrNotFound data is not exist
	ErrNotFound = gorm.ErrRecordNotFound
)

var _ Repository = (*repository)(nil)

// Repository 定义用户仓库接口
type Repository interface {
}

// repository mysql struct
type repository struct {
	orm    *gorm.DB
	db     *sql.DB
	tracer trace.Tracer
}

// New new a repository and return
func New(db *gorm.DB) Repository {
	return &repository{
		orm:    db,
		tracer: otel.Tracer("repository"),
	}
}

// Close release mysql connection
func (d *repository) Close() {

}
