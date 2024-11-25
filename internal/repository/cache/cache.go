package cache

import (
	"context"
)

type CodeCache interface {
	Set(ctx context.Context, key string, code string) error
	Verify(ctx context.Context, key, code string) (bool, error)
}
