package sms

import "context"

type Service interface {
	// Send
	// templateId 短信模板id
	// args 短信模板参数
	// phoneNums 要发送到的手机号
	Send(ctx context.Context, templateId string, args []string, phoneNums ...string) error
}
