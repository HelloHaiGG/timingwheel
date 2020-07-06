package handles

import (
	"HelloMyWorld/common"
	"HelloMyWorld/common/iredis"
	"HelloMyWorld/config"
	"HelloMyWorld/gateway"
	"HelloMyWorld/gateway/app/models"
	"HelloMyWorld/utils"
	"fmt"
	"github.com/gin-gonic/gin"
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
		"image_code", time.Duration(config.APPConfig.VerifyCode.Expiration)).Result()
	if err != nil {
		gateway.WrapErrResponse(ctx, common.OpsRedisErr, common.ServerErrDesc)
		//TODO 日志
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
		gateway.WrapErrResponse(ctx, common.ParamsErr, common.ParamsErrDesc)
	}
	if ok := utils.VerifyPhone(params.Phone); !ok {
		gateway.WrapErrResponse(ctx, common.InvalidPhoneNumberErr, common.InvalidPhoneNumberDesc)
	}
	_, err = iredis.RedisCli.Get(fmt.Sprintf(config.APPConfig.VerifyCode.Key, params.ReqId, params.ReqId)).Result()
	if err != nil {
		gateway.WrapErrResponse(ctx, common.VerifyCodeExpirationErr, common.VerifyCodeExpirationDesc)
	}
	if ok := utils.CheckCode(params.ReqId, params.Code); !ok {
		gateway.WrapErrResponse(ctx, common.InvalidCodeErr, common.InvalidCodeDesc)
	}
	//TODO kafka发送短信
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
