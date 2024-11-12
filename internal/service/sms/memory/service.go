package memory

import (
	"context"
	"fmt"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (s Service) Send(ctx context.Context, templateId string, args []string, phoneNums ...string) error {
	fmt.Println("验证码: ", args)
	return nil
}
