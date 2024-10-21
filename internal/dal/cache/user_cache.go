package cache

//go:generate mockgen -source=internal/cache/user_cache.go -destination=internal/mock/user_cache_mock.go  -package mock

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-eagle/eagle-layout/internal/dal/db/model"
	"github.com/go-eagle/eagle/pkg/cache"
	"github.com/go-eagle/eagle/pkg/encoding"
	"github.com/go-eagle/eagle/pkg/log"
	"github.com/redis/go-redis/v9"
)

const (
	// PrefixUserCacheKey cache prefix
	PrefixUserCacheKey = "user:%d"
)

// UserCache define cache interface
type UserCache interface {
	SetUserCache(ctx context.Context, id int64, data *model.UserInfoModel, duration time.Duration) error
	GetUserCache(ctx context.Context, id int64) (data *model.UserInfoModel, err error)
	MultiGetUserCache(ctx context.Context, ids []int64) (map[string]*model.UserInfoModel, error)
	MultiSetUserCache(ctx context.Context, data []*model.UserInfoModel, duration time.Duration) error
	DelUserCache(ctx context.Context, id int64) error
	SetCacheWithNotFound(ctx context.Context, id int64) error
}

// userCache define cache struct
type userCache struct {
	cache cache.Cache
}

// NewUserCache new a cache
func NewUserCache(rdb *redis.Client) UserCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""
	return &userCache{
		cache: cache.NewRedisCache(rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &model.UserInfoModel{}
		}),
	}
}

// GetUserCacheKey get cache key
func (c *userCache) GetUserCacheKey(id int64) string {
	return fmt.Sprintf(PrefixUserCacheKey, id)
}

// SetUserCache write to cache
func (c *userCache) SetUserCache(ctx context.Context, id int64, data *model.UserInfoModel, duration time.Duration) error {
	if data == nil || id == 0 {
		return nil
	}
	cacheKey := c.GetUserCacheKey(id)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// GetUserCache get from cache
func (c *userCache) GetUserCache(ctx context.Context, id int64) (data *model.UserInfoModel, err error) {
	cacheKey := c.GetUserCacheKey(id)
	err = c.cache.Get(ctx, cacheKey, &data)
	if err != nil && !errors.Is(err, redis.Nil) {
		log.WithContext(ctx).Warnf("get err from redis, err: %v", err)
		return nil, err
	}
	return data, nil
}

// MultiGetUserCache batch get cache
func (c *userCache) MultiGetUserCache(ctx context.Context, ids []int64) (map[string]*model.UserInfoModel, error) {
	var keys []string
	for _, v := range ids {
		cacheKey := c.GetUserCacheKey(v)
		keys = append(keys, cacheKey)
	}

	// NOTE: 需要在这里make实例化，如果在返回参数里直接定义会报 nil map
	retMap := make(map[string]*model.UserInfoModel)
	err := c.cache.MultiGet(ctx, keys, retMap)
	if err != nil {
		return nil, err
	}
	return retMap, nil
}

// MultiSetUserCache batch set cache
func (c *userCache) MultiSetUserCache(ctx context.Context, data []*model.UserInfoModel, duration time.Duration) error {
	valMap := make(map[string]interface{})
	for _, v := range data {
		cacheKey := c.GetUserCacheKey(v.ID)
		valMap[cacheKey] = v
	}

	err := c.cache.MultiSet(ctx, valMap, duration)
	if err != nil {
		return err
	}
	return nil
}

// DelUserCache delete cache
func (c *userCache) DelUserCache(ctx context.Context, id int64) error {
	cacheKey := c.GetUserCacheKey(id)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetCacheWithNotFound set empty cache
func (c *userCache) SetCacheWithNotFound(ctx context.Context, id int64) error {
	cacheKey := c.GetUserCacheKey(id)
	err := c.cache.SetCacheWithNotFound(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}
