package common

const (
	//-----------------------------------------网关
	RequiredErrCode = 40010 //必传参数错误
	LimiterCode     = 40011 //限流
	LimiterDesc     = "网络拥堵,请稍后再试"
	SignErrCode     = 40012 //签名错误

	//统一描述
	ServerErrDesc  = "系统错误"
	RequestErrDesc = "请求错误"
)
