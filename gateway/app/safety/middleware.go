package safety

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

var userLimiter *Limiter

func init() {
	userLimiter = Init(100, 5).Start()
}

//用户端网关限流器
func UserLimit(c *gin.Context) {
	if !userLimiter.Pass() {
		c.Next()
	} else {
		//跳转或者返回提示
		c.Redirect(http.StatusTooManyRequests,"/")
	}
}
