package gateway

import "github.com/gin-gonic/gin"

var ReqManager = make(map[string][]*WrapReq)

type WrapReq struct {
	Desc   string                 //接口描述
	Method string                 //接口请求方式
	Router string                 //接口路径
	Handle func(ctx *gin.Context) //处理方法
}

func StepForward(base string, routers []*WrapReq) {
	if _, ok := ReqManager[base]; !ok {
		ReqManager[base] = make([]*WrapReq, 0)
		ReqManager[base] = append(ReqManager[base], routers...)
	} else {
		ReqManager[base] = append(ReqManager[base], routers...)
	}
}

func Loading(engine *gin.Engine) {
	for base, reqs := range ReqManager {
		for _, req := range reqs {
			engine.Handle(req.Method, base+req.Router, req.Handle)
		}
	}
}
