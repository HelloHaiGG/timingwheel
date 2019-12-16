package safety

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

var userLimiter *Limiter

func init() {
	userLimiter = Init(100, 5).Start()
	fmt.Println("limiter init...")
}

//用户端网关限流器
func UserLimit(c *gin.Context) {
	if userLimiter.Pass() {
		c.Next()
	} else {
		//跳转或者返回提示
		c.JSON(http.StatusTooManyRequests,gin.H{
			"code":"429",
			"msg":"操作频繁,请稍后再试!",
		})
		c.Abort()
	}
}
