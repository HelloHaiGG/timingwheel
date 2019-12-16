package router

import (
	"HelloMyWorld/gateway/app/safety"
	"github.com/gin-gonic/gin"
	"net/http"
)

func UserRouter() http.Handler {
	router := gin.Default()

	//中间件
	router.Use(safety.UserLimit)

	return router
}
