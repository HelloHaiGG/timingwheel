package safety

import (
	"HelloMyWorld/common"
	"HelloMyWorld/common/consts"
	"HelloMyWorld/gateway"
	"HelloMyWorld/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
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

	if gateway.IsSafeApi(c.Request.Method, c.Request.URL.Path) {
		c.Next()
		return
	}
	ts := c.Request.Header.Get(consts.TimeStamp)
	version := c.Request.Header.Get(consts.Version)
	sign := c.Request.Header.Get(consts.Sign)
	required = fmt.Sprintf("Version=%sTs=%s", version, ts)

	//判断排序方式
	timeStamp, _ := strconv.Atoi(ts)
	sortWay := timeStamp % 2 //0：正序 1：倒序
	params = spliceParams(c.Request, cast.ToInt32(sortWay))
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

//判断必须参数
func verifyRequiredParams(header http.Header) bool {
	if ts, err := cast.ToInt64E(header.Get(consts.TimeStamp)); err != nil || ts == 0 {
		return false
	}
	if version := header.Get(consts.Version); version == "" {
		return false
	}
	if sign := header.Get(consts.Sign); sign == "" {
		return false
	}
	return true
}

//组合参数字符串
func spliceParams(r *http.Request, t int32) string {
	var params string
	switch r.Method {
	case "GET", "HEAD":
		params = r.URL.Query().Encode()
	case "POST", "PUT":
		params = r.PostForm.Encode()
	case "DELETE":
		//对于delete请求,少量参数的情况下可以使用,参数较多还是使用post请求
		params = r.PostForm.Encode()
	default:
		//TODO
	}

	if t == 0 {
		return params
	} else {
		p := strings.Split(params, "&")
		return strings.Join(utils.SortStr(p, -1), "&")
	}
}
