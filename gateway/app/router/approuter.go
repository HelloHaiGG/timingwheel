package router

import (
	"HelloMyWorld/gateway"
	"HelloMyWorld/gateway/app/safety"
	"github.com/gin-gonic/gin"
	"net/http"
)

func init() {
	gateway.StepForward("/jx/app", []*gateway.WrapReq{
		{Desc: "", Method: "", Router: "", Handle: nil,},
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

	router.GET("/api", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"msg": context.Param("a"),
		})
	})
	router.GET("/", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"msg": "重定向",
		})
	})

	return router
}
