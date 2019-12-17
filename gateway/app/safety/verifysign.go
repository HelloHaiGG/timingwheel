package safety

import (
	"HelloMyWorld/common"
	"HelloMyWorld/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

var (
	required string //必传参数
	params   string //参数
	target   string //需要进行md5计算的目标字符串
)

//功能:防止恶意截取接口,篡改参数
//签名策略
//前端参数传递必传参数有 TS,Sign,Version
//通过TS对2取余,判断参数排序方式 0：正序 1：倒序
//将参数组合成字符串,生成sign
//判断前端穿过来的sign与生成的sign是否一样

func VerifySign(c *gin.Context) {
	//必传参数
	Ts := c.Request.Header.Get("Ts")
	Version := c.Request.Header.Get("JX-Version")
	sign := c.Request.Header.Get("Sign")
	required = fmt.Sprintf("Version=%sTs=%s", Version, Ts)
	//判断排序方式
	timeStamp, _ := strconv.Atoi(Ts)
	sortWay := timeStamp % 2 //0：正序 1：倒序

	switch c.Request.Method {
	case "GET", "HEAD":
		if sortWay == 0 {
			params = c.Request.URL.Query().Encode()
		} else {
			p := strings.Split(c.Request.URL.Query().Encode(), "&")
			params = strings.Join(utils.SortStr(p, -1), "&")
		}
	case "POST", "PUT":
	case "DELETE":
	default:
		//TODO
	}
	target = fmt.Sprintf("%s&%s", required, params)
	//将目标字符串进行md5算法
	s := utils.GetMd5WithStr(target)
	//比较与前端传的sign 是否一致
	if strings.Compare(s, sign) != 0 {
		c.JSON(http.StatusForbidden, gin.H{
			"code": common.SignErrCode,
			"msg":  common.RequiredErrCode,
		})
		c.Abort()
	}
	c.Next()
}
