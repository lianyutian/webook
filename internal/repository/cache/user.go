package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
	"webook/internal/domain"
)

var (
	ErrKeyNotExits = redis.Nil
)

type UserCache struct {
	client     redis.Cmdable
	expiration time.Duration
}

// NewUserCache
// A 用到了 B, B 一定是接口 ==》保证面向接口
// A 用到了 B, B 一定是 A 的字段 ==》 规避报变量、包方法
// A 用到了 B, A 绝不初始化 B, 而是由外部注入 ==》 保持依赖注入和依赖反转
func NewUserCache(client redis.Cmdable) *UserCache {
	return &UserCache{client: client}
}

// Get
// 如果返回不是err，一定存在数据
func (cache *UserCache) Get(ctx context.Context, id int64) (domain.User, error) {
	key := cache.key(id)
	val, err := cache.client.Get(ctx, key).Bytes()
	if err != nil {
		return domain.User{}, err
	}
	err = json.Unmarshal(val, &domain.User{})
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{}, nil
}

func (cache *UserCache) Set(ctx context.Context, id int64, user domain.User) error {
	val, err := json.Marshal(user)
	if err != nil {
		return err
	}
	key := cache.key(id)
	return cache.client.Set(ctx, key, val, cache.expiration).Err()
}

func (cache *UserCache) key(id int64) string {
	return fmt.Sprintf("user:info:%d", id)
}
