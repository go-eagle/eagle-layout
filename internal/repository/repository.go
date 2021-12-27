package repository

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
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
	db     *gorm.DB
	tracer trace.Tracer
}

// New new a repository and return
func New(db *gorm.DB) Repository {
	return &repository{
		db:     db,
		tracer: otel.Tracer("repository"),
	}
}

// Close release mysql connection
func (d *repository) Close() {

}
