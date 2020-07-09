package gateway

import (
	"HelloMyWorld/common"
	"HelloMyWorld/common/ilogger"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Result struct {
	Code int         `json:"Code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"Data"`
}

func SucResponse(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, &Result{
		Code: common.Success,
		Msg:  common.SuccessDesc,
	})
}

func WrapResponse(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, &Result{
		Code: common.Success,
		Msg:  common.SuccessDesc,
		Data: data,
	})
}

func WrapErrRequest(ctx *gin.Context, httpCode int) {
	ctx.JSON(httpCode, nil)
	ctx.Abort()
}

func WrapErrResponse(ctx *gin.Context, err error, code int, msg string) {
	ilogger.Ins.Error(err)
	ctx.JSON(http.StatusOK, &Result{
		Code: code,
		Msg:  msg,
		Data: nil,
	})
	ctx.Abort()
}
