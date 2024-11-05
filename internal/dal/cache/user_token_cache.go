package cache

//go:generate mockgen -source=internal/cache/user_token_cache.go -destination=internal/mock/user_token_cache_mock.go  -package mock

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-eagle/eagle/pkg/cache"
	"github.com/go-eagle/eagle/pkg/encoding"
	"github.com/go-eagle/eagle/pkg/log"
	"github.com/go-eagle/eagle/pkg/redis"
)

const (
	// PrefixUserTokenCacheKey cache prefix
	PrefixUserTokenCacheKey = "user:token:%d"
	UserTokenExpireTime     = 24 * time.Hour * 30
)

// UserToken define cache interface
type UserTokenCache interface {
	SetUserTokenCache(ctx context.Context, id int64, token string, duration time.Duration) error
	GetUserTokenCache(ctx context.Context, id int64) (token string, err error)
	DelUserTokenCache(ctx context.Context, id int64) error
}

// userTokenCache define cache struct
type userTokenCache struct {
	cache cache.Cache
}

// NewUserTokenCache new a cache
func NewUserTokenCache() UserTokenCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""
	return &userTokenCache{
		cache: cache.NewRedisCache(redis.RedisClient, cachePrefix, jsonEncoding, func() interface{} {
			return ""
		}),
	}
}

// GetUserTokenCacheKey get cache key
func (c *userTokenCache) GetUserTokenCacheKey(id int64) string {
	return fmt.Sprintf(PrefixUserTokenCacheKey, id)
}

// SetUserTokenCache write to cache
func (c *userTokenCache) SetUserTokenCache(ctx context.Context, id int64, token string, duration time.Duration) error {
	if len(token) == 0 || id == 0 {
		return nil
	}
	cacheKey := c.GetUserTokenCacheKey(id)
	err := c.cache.Set(ctx, cacheKey, &token, duration)
	if err != nil {
		return err
	}
	return nil
}

// GetUserTokenCache get from cache
func (c *userTokenCache) GetUserTokenCache(ctx context.Context, id int64) (token string, err error) {
	cacheKey := c.GetUserTokenCacheKey(id)
	err = c.cache.Get(ctx, cacheKey, &token)
	if err != nil && !errors.Is(err, redis.ErrRedisNotFound) {
		log.WithContext(ctx).Warnf("get err from redis, err: %+v", err)
		return "", err
	}
	return token, nil
}

// DelUserTokenCache delete cache
func (c *userTokenCache) DelUserTokenCache(ctx context.Context, id int64) error {
	cacheKey := c.GetUserTokenCacheKey(id)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}
