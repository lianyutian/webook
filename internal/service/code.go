package service

import (
	"context"
	"fmt"
	"math/rand"
	"webook/internal/repository"
	"webook/internal/service/sms"
)

var (
	ErrCodeSendTooMany        = repository.ErrCodeSendTooMany
	ErrCodeVerifyTooManyTimes = repository.ErrCodeVerifyTooManyTimes
)

type CodeService struct {
	rep    *repository.CodeRepository
	smsSvc sms.Service
}

func NewCodeService(rep *repository.CodeRepository, smsSvc sms.Service) *CodeService {
	return &CodeService{
		rep:    rep,
		smsSvc: smsSvc,
	}
}

func (svc *CodeService) Send(ctx context.Context, templateId, biz string, phoneNum string) error {
	// 生成 code 码
	code := svc.generateCode()
	// code 存储 redis
	err := svc.rep.Store(ctx, biz, phoneNum, code)
	if err != nil {
		return err
	}
	// 发送短信
	err = svc.smsSvc.Send(ctx, templateId, []string{code}, phoneNum)
	if err != nil {
		// 如果前面写入 redis 成功, 这里失败。用户实际上是没有收到短信的
		// 能不能删除 redis 里的验证码？
		// 不能，因为这个错误不能确定原因，可能是短信已经发送出去了但是接收返回信息时超时了
		// TODO 所以这里需要有重试机制
	}
	return nil
}

func (svc *CodeService) Verify(ctx context.Context, biz string, phoneNum string, code string) (bool, error) {
	return svc.rep.Verify(ctx, biz, phoneNum, code)
}

func (svc *CodeService) generateCode() string {
	// 生成 [0-1000000) 之间的数
	num := rand.Intn(1000000)
	// 不够 6 位将前位补 0
	return fmt.Sprintf("%06d", num)
}
