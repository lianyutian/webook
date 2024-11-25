package repository

import (
	"context"
	"webook/internal/repository/cache"
)

var (
	ErrCodeSendTooMany        = cache.ErrCodeSendTooMany
	ErrCodeVerifyTooManyTimes = cache.ErrCodeVerifyTooManyTimes
)

type CodeRepository struct {
	cache *cache.CodeRedisCache
}

func NewCodeRepository(cache *cache.CodeRedisCache) *CodeRepository {
	return &CodeRepository{
		cache: cache,
	}
}

func (svc *CodeRepository) Store(ctx context.Context, biz, phoneNum, code string) error {
	return svc.cache.Set(ctx, biz, phoneNum, code)
}

func (svc *CodeRepository) Verify(ctx context.Context, biz, phoneNum, code string) (bool, error) {
	return svc.cache.Verify(ctx, biz, phoneNum, code)
}
