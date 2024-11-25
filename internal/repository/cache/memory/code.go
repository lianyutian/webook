package memory

import (
	"context"
	"github.com/allegro/bigcache/v3"
)

type CodeCache struct {
	bigCache *bigcache.BigCache
}

func (cache *CodeCache) Set(ctx context.Context, key string, code string) error {
	err := cache.bigCache.Set(key, []byte(code))
	if err != nil {
		return err
	}
	return nil
}

func (cache *CodeCache) Verify(ctx context.Context, key, code string) (bool, error) {
	val, err := cache.bigCache.Get(key)
	if err != nil {
		return false, err
	}
	if string(val) == code {
		return true, nil
	}
	return false, nil
}
