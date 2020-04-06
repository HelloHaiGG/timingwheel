package safety

import (
	"HelloMyWorld/common"
	"github.com/gin-gonic/gin"
	"net/http"
)

var limiter *Limiter

func init() {
	limiter = Init(100, 5).Start()
}

//用户端网关限流器
func Limit(c *gin.Context) {
	if limiter.Pass() {
		c.Next()
	} else {
		//跳转或者返回提示
		c.JSON(http.StatusTooManyRequests, gin.H{
			"code": common.LimiterCode,
			"msg":  common.LimiterDesc,
		})
		c.Abort()
	}
}
