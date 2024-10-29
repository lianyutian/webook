package ali

//
//import (
//	"context"
//	dysmsapi "github.com/alibabacloud-go/dysmsapi-20170525/v2/client"
//)
//
//type Service struct {
//	client *dysmsapi.Client
//	appId  *string
//	sign   *string
//}
//
//func NewService(client *dysmsapi.Client, appId, sign *string) *Service {
//	return &Service{
//		client: client,
//		appId:  appId,
//		sign:   sign,
//	}
//}
//
//func (s *Service) Send(ctx context.Context, templateId string, args []string, phoneNums ...string) error {
//	// 1.发送短信请求
//	sendReq := &dysmsapi.SendSmsRequest{
//		PhoneNumbers:  s.toStringPtrSlice(phoneNums),
//		SignName:      s.sign,
//		TemplateCode:  &templateId,
//		TemplateParam: args[3],
//	}
//	// 2.发送短信
//	sendResp, _err := s.client.SendSms(sendReq)
//	if _err != nil {
//		return _err
//	}
//}
//
//func (s *Service) toStringPtrSlice(src []string) *[]string {
//	return slice.Map[string, *string](src, func(idx int, value string) *string {
//		return &src
//	})
//}
