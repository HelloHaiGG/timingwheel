package safety

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

/**
必传参数验证
*/

func RequiredParams(c *gin.Context) {
	if len(c.Request.Header.Get("Ts")) == 0 {
		c.JSON(http.StatusForbidden, gin.H{
			"code": "10010",
			"msg":  "错误请求",
		})
		c.Abort()
	}
	if len(c.Request.Header.Get("JX-Version")) == 0 {
		c.JSON(http.StatusForbidden, gin.H{
			"code": "10010",
			"msg":  "错误请求",
		})
		c.Abort()
	}
	if len(c.Request.Header.Get("Sign")) == 0 {
		c.JSON(http.StatusForbidden, gin.H{
			"code": "10010",
			"msg":  "错误请求",
		})
		c.Abort()
	}
}
