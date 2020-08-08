package tencentcloud

import (
	common2 "HelloMyWorld/common"
	"fmt"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20190711"
	"log"
	"strings"
)

var SmsClient *sms.Client
var credential *common.Credential

func Init(id, key string) {
	credential = common.NewCredential(id, key)
	cpf := profile.NewClientProfile()
	client, err := sms.NewClient(credential, "ap-guangzhou", cpf)
	if err != nil {
		log.Fatal("Tencent cloud send sms client err:", err)
	}
	fmt.Println("Tencent cloud connect.")
	SmsClient = client
}

func SendSms(phones []string, params []string) error {
	if len(phones) > 200 || len(phones) == 0 {
		return nil
	} else {
		for i, phone := range phones {
			phones[i] = strings.Join([]string{"+86", phone}, "")
		}
	}
	if SmsClient == nil {
		//TODO
		fmt.Println("tencent client is nil.")
		return nil
	}
	req := sms.NewSendSmsRequest()
	req.SmsSdkAppid = common.StringPtr(common2.TenSMSAppId)
	req.Sign = common.StringPtr("禾健网络科技")
	req.PhoneNumberSet = common.StringPtrs(phones)
	req.TemplateID = common.StringPtr("673145")
	req.TemplateParamSet = common.StringPtrs(params)
	rsp, err := SmsClient.SendSms(req)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		//sdk 异常
		return err
	}
	if err != nil {
		//TODO 发送失败
		fmt.Println("tencent send sms err:", err)
	}

	return nil
}
