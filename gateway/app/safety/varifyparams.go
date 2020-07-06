package safety

import (
	"HelloMyWorld/common"
	"HelloMyWorld/gateway"
	"github.com/gin-gonic/gin"
	"net/http"
)

/**
必传参数验证
*/

func RequiredParams(c *gin.Context) {
	if gateway.IsSafeApi(c.Request.Method, c.Request.URL.Path) {
		c.Next()
		return
	}
	if len(c.Request.Header.Get("Ts")) == 0 {
		c.JSON(http.StatusForbidden, gin.H{
			"code": common.RequiredErrCode,
			"msg":  common.RequestErrDesc,
		})
		c.Abort()
		return
	}
	if len(c.Request.Header.Get("JX-Version")) == 0 {
		c.JSON(http.StatusForbidden, gin.H{
			"code": common.RequiredErrCode,
			"msg":  common.RequestErrDesc,
		})
		c.Abort()
		return
	}
	if len(c.Request.Header.Get("Sign")) == 0 {
		c.JSON(http.StatusForbidden, gin.H{
			"code": common.RequiredErrCode,
			"msg":  common.RequestErrDesc,
		})
		c.Abort()
		return
	}
}
