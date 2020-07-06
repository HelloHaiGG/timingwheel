package router

import (
	"HelloMyWorld/gateway"
	"HelloMyWorld/gateway/app/safety"
	"github.com/gin-gonic/gin"
	"net/http"
)

func init() {
	gateway.StepForward("/jx/app", []*gateway.WrapReq{
		{Desc: "服务测试", Method: "GET", Router: "/health", Handle: ServerStatus, Safe: true},
		{Desc: "服务测试", Method: "GET", Router: "", Handle: ServerStatus, Safe: true},
	})
	gateway.StepForward("/", []*gateway.WrapReq{
		{Desc: "服务测试", Method: "GET", Router: "", Handle: ServerStatus, Safe: true},
	})
}

func APPRouter() *gin.Engine {
	router := gin.Default()

	//中间件
	router.Use(
		//限流器
		safety.Limit,
		//必要参数验证
		safety.RequiredParams,
		//签名验证
		safety.VerifySign,
	)
	//加载路由
	gateway.Loading(router)

	return router
}

func ServerStatus(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "Hello,Jx.")
}
