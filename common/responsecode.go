package common

const (
	Success = 10000

	//-----------------------------------------网关
	RequiredErrCode = 400100 //必传参数错误
	LimiterCode     = 400101 //限流
	LimiterDesc     = "网络拥堵,请稍后再试"
	SignErrCode     = 400102 //签名错误

	//-----------------------------------------统一错误
	ParamsErr = 500001

	//-----------------------------------------Redis
	OpsRedisErr = 637900

	//-----------------------------------------定义错误
	InvalidPhoneNumberErr   = 500100
	InvalidCodeErr          = 500101
	VerifyCodeExpirationErr = 500102
)

const (
	SuccessDesc              = "成功"
	InvalidPhoneNumberDesc   = "无效的手机号"
	ServerErrDesc            = "系统错误"
	ParamsErrDesc            = "参数错误"
	RequestErrDesc           = "请求错误"
	InvalidCodeDesc          = "验证码无效"
	VerifyCodeExpirationDesc = "验证码已过期"
)
