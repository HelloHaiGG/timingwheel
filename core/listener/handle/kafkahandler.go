package handle

import (
	sms "HelloMyWorld/common/entity/kafka"
	"HelloMyWorld/common/ilogger"
	"HelloMyWorld/common/tencentcloud"
	"fmt"
	"github.com/gogo/protobuf/proto"
)

func Handler(msg proto.Message) {
	switch msg.(type) {
	case *sms.SendSms:
		entity := msg.(*sms.SendSms)
		SendSms(entity)
	}
}

func SendSms(sms *sms.SendSms) {
	if err := tencentcloud.SendSms([]string{sms.Phone}, []string{sms.Code, sms.EffeTime}); err != nil {
		ilogger.Ins.Error(err)
		return
	}
	ilogger.Ins.Info(fmt.Sprintf("Send sms code to %s:%s success.", sms.Phone, sms.Code))
}
