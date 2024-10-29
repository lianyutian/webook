package cache

import (
	"context"
	_ "embed"
	"errors"
	"github.com/redis/go-redis/v9"
)

var (
	ErrCodeSendTooMany        = errors.New("验证码发送太频繁")
	ErrCodeVerifyTooManyTimes = errors.New("验证次数太多")
	ErrUnknownForCode         = errors.New("未知错误")
)

// 编译器会在编译的时候，把 set_code.lua 的代码放入 luaSetCode 变量
//
//go:embed lua/set_code.lua
var luaSetCode string

//go:embed lua/verify_code.lua
var luaVerifyCode string

type CodeCache struct {
	client redis.Cmdable
}

func NewCodeCache(client redis.Cmdable) *CodeCache {
	return &CodeCache{client: client}
}

// Set
// 执行 lua 脚本
func (cache *CodeCache) Set(ctx context.Context, biz, phoneNum, code string) error {
	res, err := cache.client.Eval(ctx, luaSetCode, []string{biz, phoneNum}, code).Int()
	if err != nil {
		return err
	}
	switch res {
	case 0:
		return nil
	case -1: // 发送频繁
		return ErrCodeSendTooMany
	//case -2: // 系统错误
	default:
		return errors.New("系统错误")
	}
}

func (cache *CodeCache) Verify(ctx context.Context, biz, phoneNum, code string) (bool, error) {
	res, err := cache.client.Eval(ctx, luaVerifyCode, []string{biz, phoneNum}, code).Int()
	if err != nil {
		return false, err
	}
	switch res {
	case 0:
		return true, nil
	case -1:
		// TODO 如果频繁出现这个报错就需要告警
		return false, ErrCodeVerifyTooManyTimes
	case -2:
		return false, nil
	default:
		return false, ErrUnknownForCode
	}
}
