package repository

//go:generate mockgen -source=user_repo.go -destination=../../internal/mocks/user_repo_mock.go  -package mocks

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	"github.com/go-eagle/eagle-layout/internal/cache"
	"github.com/go-eagle/eagle-layout/internal/dal/model"
)

var (
	_tableUserName   = (&model.UserModel{}).TableName()
	_getUserSQL      = "SELECT * FROM %s WHERE id = ?"
	_batchGetUserSQL = "SELECT * FROM %s WHERE id IN (%s)"
)

var _ UserRepo = (*userRepo)(nil)

// UserRepo define a repo interface
type UserRepo interface {
	CreateUser(ctx context.Context, data *model.UserModel) (id int64, err error)
	UpdateUser(ctx context.Context, id int64, data *model.UserModel) error
	GetUser(ctx context.Context, id int64) (ret *model.UserModel, err error)
	BatchGetUser(ctx context.Context, ids []int64) (ret []*model.UserModel, err error)
}

type userRepo struct {
	db     *model.DBClient
	tracer trace.Tracer
	cache  cache.UserCache
}

// NewUser new a repository and return
func NewUserRepo(db *model.DBClient, cache cache.UserCache) UserRepo {
	return &userRepo{
		db:     db,
		tracer: otel.Tracer("user"),
		cache:  cache,
	}
}

// CreateUser create a item
func (r *userRepo) CreateUser(ctx context.Context, data *model.UserModel) (id int64, err error) {
	err = r.db.GetDB().WithContext(ctx).Create(&data).Error
	if err != nil {
		return 0, errors.Wrap(err, "[repo] create User err")
	}

	return data.ID, nil
}

// UpdateUser update item
func (r *userRepo) UpdateUser(ctx context.Context, id int64, data *model.UserModel) error {
	item, err := r.GetUser(ctx, id)
	if err != nil {
		return errors.Wrapf(err, "[repo] update User err: %v", err)
	}
	err = r.db.GetDB().Model(&item).Updates(data).Error
	if err != nil {
		return err
	}
	// delete cache
	_ = r.cache.DelUserCache(ctx, id)
	return nil
}

// GetUser get a record
func (r *userRepo) GetUser(ctx context.Context, id int64) (ret *model.UserModel, err error) {
	// read cache
	item, err := r.cache.GetUserCache(ctx, id)
	if err != nil {
		return nil, err
	}
	if item != nil {
		return item, nil
	}
	// read db
	data := new(model.UserModel)
	err = r.db.GetDB().WithContext(ctx).Raw(fmt.Sprintf(_getUserSQL, _tableUserName), id).Scan(&data).Error
	if err != nil {
		return
	}
	// write cache
	if data.ID > 0 {
		err = r.cache.SetUserCache(ctx, id, data, 5*time.Minute)
		if err != nil {
			return nil, err
		}
	}
	return data, nil
}

// BatchGetUser batch get items
func (r *userRepo) BatchGetUser(ctx context.Context, ids []int64) (ret []*model.UserModel, err error) {
	// read cache
	idsStr := cast.ToStringSlice(ids)
	itemMap, err := r.cache.MultiGetUserCache(ctx, ids)
	if err != nil {
		return nil, err
	}
	var missedID []int64
	for _, v := range ids {
		item, ok := itemMap[cast.ToString(v)]
		if !ok {
			missedID = append(missedID, v)
			continue
		}
		ret = append(ret, item)
	}
	// get missed data
	if len(missedID) > 0 {
		var missedData []*model.UserModel
		_sql := fmt.Sprintf(_batchGetUserSQL, _tableUserName, strings.Join(idsStr, ","))
		err = r.db.GetDB().WithContext(ctx).Raw(_sql).Scan(&missedData).Error
		if err != nil {
			// you can degrade to ignore error
			return nil, err
		}
		if len(missedData) > 0 {
			ret = append(ret, missedData...)
			err = r.cache.MultiSetUserCache(ctx, missedData, 5*time.Minute)
			if err != nil {
				// you can degrade to ignore error
				return nil, err
			}
		}
	}
	return ret, nil
}
