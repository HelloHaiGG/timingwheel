package router

import (
	"HelloMyWorld/gateway/app/safety"
	"github.com/gin-gonic/gin"
	"net/http"
)

func UserRouter() http.Handler {
	router := gin.Default()

	//中间件
	router.Use(
		safety.UserLimit,
		//必要参数验证
		safety.RequiredParams,
		//签名验证
		safety.VerifySign,
	)

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
