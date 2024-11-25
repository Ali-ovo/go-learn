package serviceSms

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	srvSms "shop/app/shop/api/internal_api/service/sms"
	"shop/pkg/options"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
)

type SmsDTO struct{}

type smsService struct {
	smsOpts *options.SmsOptions
}

func GenerateSmsCode(witdh int) string {
	//生成width长度的短信验证码

	numeric := [10]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	r := len(numeric)
	rand.NewSource(time.Now().UnixNano())

	var sb strings.Builder
	for i := 0; i < witdh; i++ {
		fmt.Fprintf(&sb, "%d", numeric[rand.Intn(r)])
	}
	return sb.String()
}

// SendSms
//
//	@Description: 发送短信验证码
//	@param ctx
//	@param mobile: 手机号码
//	@param tpc: template code 消息模板编号
//	@param tp: template param 消息参数
//	@return error
func (s *smsService) SendSms(ctx context.Context, mobile string, tp string) error {
	client, err := dysmsapi.NewClientWithAccessKey("cn-beijing", s.smsOpts.APIKey, s.smsOpts.APISecret)
	if err != nil {
		return err
	}
	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Scheme = "https" // https | http
	request.Domain = "dysmsapi.aliyuncs.com"
	request.Version = "2017-05-25"
	request.ApiName = "SendSms"
	request.QueryParams["RegionId"] = "cn-beijing"
	request.QueryParams["PhoneNumbers"] = mobile            //手机号
	request.QueryParams["SignName"] = s.smsOpts.APISignName //阿里云验证过的项目名 自己设置
	request.QueryParams["TemplateCode"] = s.smsOpts.APICode //阿里云的短信模板号 自己设置
	request.QueryParams["TemplateParam"] = tp               //短信模板中的验证码内容 自己生成   之前试过直接返回，但是失败，加上code成功。
	response, err := client.ProcessCommonRequest(request)
	if err != nil {
		return err
	}

	isSendBytes := response.GetHttpContentBytes()
	isSend := make(map[string]interface{})
	err = json.Unmarshal(isSendBytes, &isSend)
	if err != nil {
		return err
	}

	if isSend["Message"] != "OK" {
		return fmt.Errorf("发送短信验证码失败")
	}
	return nil
}

func NewSmsService(smsOpts *options.SmsOptions) srvSms.SmsSrv {
	return &smsService{smsOpts: smsOpts}
}
