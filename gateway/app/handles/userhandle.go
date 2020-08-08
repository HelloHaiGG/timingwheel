package handles

import (
	"HelloMyWorld/common"
	sms "HelloMyWorld/common/entity/kafka"
	"HelloMyWorld/common/ikafka"
	"HelloMyWorld/common/ilogger"
	"HelloMyWorld/common/iredis"
	"HelloMyWorld/config"
	"HelloMyWorld/gateway"
	"HelloMyWorld/gateway/app/models"
	"HelloMyWorld/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gogo/protobuf/proto"
	"time"
)

func init() {
	gateway.StepForward("/jx/app", []*gateway.WrapReq{
		{Desc: "用户获取验证码", Method: "GET", Router: "/code", Handle: GetCode},
		{Desc: "发送短信验证码", Method: "POST", Router: "/sms", Handle: SendSMS},
		{Desc: "短信验证码验证", Method: "POST", Router: "/sms/check", Handle: CheckSMS},
		{Desc: "用户注册", Method: "POST", Router: "/register", Handle: Register,},
	})
}

func GetCode(ctx *gin.Context) {
	reqId, image := utils.CreateImageCode()
	_, err := iredis.RedisCli.Set(
		fmt.Sprintf(config.APPConfig.VerifyCode.Key, reqId),
		"image_code", time.Duration(config.APPConfig.VerifyCode.Expiration)*time.Second).Result()
	if err != nil {
		gateway.WrapErrResponse(ctx, err, common.OpsRedisErr, common.ServerErrDesc)
		ilogger.Ins.Error(err)
		return
	}
	gateway.WrapResponse(ctx, struct {
		ID    string `json:"id"`
		Image string `json:"image"`
	}{reqId, image})
}

func SendSMS(ctx *gin.Context) {
	var params models.SendSMSParams
	err := ctx.ShouldBindJSON(&params)
	if err != nil {
		gateway.WrapErrResponse(ctx, err, common.ParamsErr, common.ParamsErrDesc)
		return
	}
	if ok := utils.VerifyPhone(params.Phone); !ok {
		gateway.WrapErrResponse(ctx, err, common.InvalidPhoneNumberErr, common.InvalidPhoneNumberDesc)
		return
	}
	//_, err = iredis.RedisCli.Get(fmt.Sprintf(config.APPConfig.VerifyCode.Key, params.ReqId)).Result()
	//if err != nil {
	//	gateway.WrapErrResponse(ctx, err, common.InvalidCodeErr, common.InvalidCodeDesc)
	//	return
	//}
	//if ok := utils.CheckCode(params.ReqId, params.Code); !ok {
	//	gateway.WrapErrResponse(ctx, nil, common.InvalidCodeErr, common.InvalidCodeDesc)
	//	return
	//}
	//验证通过,清除
	iredis.RedisCli.Del(fmt.Sprintf(config.APPConfig.VerifyCode.Key, params.ReqId))
	//生成短信验证码
	sCode, err := utils.Ins.RandStrWithSeedChar("0,1,2,3,4,5,6,7,8,9", time.Now().Unix())
	if err != nil {
		ilogger.Ins.Error(err)
	}
	ikafka.Kafka.ASyncSendMsg(&ikafka.KafkaMsg{
		Topic: common.SendSMSTopic,
		Key:   common.SendSMSKey,
		Value: proto.MarshalTextString(&sms.SendSms{Phone: params.Phone, Code: sCode, EffeTime: "5"}),
	})
	_, err = iredis.RedisCli.Set(params.Phone, sCode, 5*60*time.Second).Result()
	if err != nil {
		ilogger.Ins.Error(err)
		return
	}
	gateway.SucResponse(ctx)
}
func CheckSMS(ctx *gin.Context) {

}

func Register(ctx *gin.Context) {
	var params models.RegisterParams
	err := ctx.ShouldBindJSON(&params)
	if err != nil {

	}
}
