package tencent

import (
	"context"
	"github.com/go-playground/assert/v2"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
	"testing"
)

func TestService_Send(t *testing.T) {
	// 替换为你的 SecretId 和 SecretKey
	secretId := "YOUR_SECRET_ID"
	secretKey := "YOUR_SECRET_KEY"

	// 创建一个腾讯云 SMS 客户端
	credential := common.NewCredential(secretId, secretKey)
	cpf := profile.NewClientProfile()
	cpf.SignMethod = "HmacSHA256"
	client, _ := sms.NewClient(credential, "ap-guangzhou", cpf)

	s := NewService(client, "appId", "sign")

	type S struct {
		name string
	}

	testCases := []struct {
		name       string
		templateId string
		params     []string
		phoneNums  []string
		wantErr    error
	}{
		{
			name:       "发送验证码",
			templateId: "模板id",
			params:     []string{"123456"},
			phoneNums:  []string{"1572XX69"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := s.Send(context.Background(), tc.templateId, tc.params, tc.phoneNums...)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}
