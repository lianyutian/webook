package tencent

import (
	"context"
	"fmt"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
)

type Service struct {
	client *sms.Client
	appId  *string
	sign   *string
}

func NewService(client *sms.Client, appId string, sign string) *Service {
	return &Service{
		client: client,
		appId:  &appId,
		sign:   &sign,
	}
}

func (s *Service) Send(ctx context.Context, templateId string, args []string, phoneNums ...string) error {
	// 创建发送请求
	request := sms.NewSendSmsRequest()

	request.SetContext(ctx)
	request.SmsSdkAppId = s.appId
	request.SignName = s.sign                             // 短信签名
	request.TemplateId = common.StringPtr(templateId)     // 短信模板 ID
	request.TemplateParamSet = common.StringPtrs(args)    // 短信模板参数
	request.PhoneNumberSet = common.StringPtrs(phoneNums) // 发送到的手机号

	// 发送短信
	response, err := s.client.SendSms(request)
	if err != nil {
		return err
	}
	for _, status := range response.Response.SendStatusSet {
		if status.Code == nil || *(status.Code) != "Ok" {
			return fmt.Errorf("send code faild code: %s message: %s", *status.Code, *status.Message)
		}
	}
	return nil
}
