package repository

//go:generate mockgen -source=user_repo.go -destination=../../internal/mocks/user_repo_mock.go  -package mocks

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"

	"github.com/go-eagle/eagle-layout/internal/dal"
	"github.com/go-eagle/eagle-layout/internal/dal/cache"
	"github.com/go-eagle/eagle-layout/internal/dal/db/dao"
	"github.com/go-eagle/eagle-layout/internal/dal/db/model"
)

var _ UserRepo = (*userRepo)(nil)

// UserRepo define a repo interface
type UserRepo interface {
	CreateUser(ctx context.Context, data model.UserInfoModel) (id int64, err error)
	UpdateUser(ctx context.Context, id int64, data model.UserInfoModel) error
	GetUser(ctx context.Context, id int64) (ret *model.UserInfoModel, err error)
	BatchGetUsers(ctx context.Context, ids []int64) (ret []*model.UserInfoModel, err error)
	GetUserByUsername(ctx context.Context, username string) (ret *model.UserInfoModel, err error)
	GetUserByEmail(ctx context.Context, email string) (ret *model.UserInfoModel, err error)
	GetUserByPhone(ctx context.Context, phone string) (ret *model.UserInfoModel, err error)
}

type userRepo struct {
	db     *dal.DBClient
	tracer trace.Tracer
	cache  cache.UserCache
}

// NewUser new a repository and return
func NewUserRepo(db *dal.DBClient, cache cache.UserCache) UserRepo {
	return &userRepo{
		db:     db,
		tracer: otel.Tracer("user"),
		cache:  cache,
	}
}

// CreateUser create a item
func (r *userRepo) CreateUser(ctx context.Context, data model.UserInfoModel) (id int64, err error) {
	err = dao.UserInfoModel.WithContext(ctx).Create(&data)
	if err != nil {
		return 0, errors.Wrap(err, "[repo] create User err")
	}

	return data.ID, nil
}

// UpdateUser update item
func (r *userRepo) UpdateUser(ctx context.Context, id int64, data model.UserInfoModel) error {
	_, err := dao.UserInfoModel.WithContext(ctx).Where(dao.UserInfoModel.ID.Eq(id)).Updates(data)
	if err != nil {
		return err
	}
	// delete cache
	_ = r.cache.DelUserCache(ctx, id)
	return nil
}

// GetUser get a record
func (r *userRepo) GetUser(ctx context.Context, id int64) (ret *model.UserInfoModel, err error) {
	// read cache
	item, err := r.cache.GetUserCache(ctx, id)
	if err != nil {
		return nil, err
	}
	if item != nil {
		return item, nil
	}

	// read db
	data, err := dao.UserInfoModel.WithContext(ctx).Where(dao.UserInfoModel.ID.Eq(id)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.cache.SetCacheWithNotFound(ctx, id)
		}
		return
	}

	// write cache
	if data != nil && data.ID > 0 {
		err = r.cache.SetUserCache(ctx, id, data, 5*time.Minute)
		if err != nil {
			return nil, err
		}
	}
	return data, nil
}

// BatchGetUser batch get items
func (r *userRepo) BatchGetUsers(ctx context.Context, ids []int64) (ret []*model.UserInfoModel, err error) {
	// read cache
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
		missedData, err := dao.UserInfoModel.WithContext(ctx).Where(dao.UserInfoModel.ID.In(ids...)).Find()
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

func (r *userRepo) GetUserByUsername(ctx context.Context, username string) (ret *model.UserInfoModel, err error) {
	ret, err = dao.UserInfoModel.WithContext(ctx).GetUserByUsername(username)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *userRepo) GetUserByEmail(ctx context.Context, email string) (ret *model.UserInfoModel, err error) {
	ret, err = dao.UserInfoModel.WithContext(ctx).GetUserByEmail(email)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *userRepo) GetUserByPhone(ctx context.Context, phone string) (ret *model.UserInfoModel, err error) {
	ret, err = dao.UserInfoModel.WithContext(ctx).GetUserByPhone(phone)
	if err != nil {
		return
	}

	return ret, nil
}
