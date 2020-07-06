package gateway

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
)

var ReqManager = make(map[string][]*WrapReq)
var SafeApi = make(map[string]*WrapReq)

type WrapReq struct {
	Desc   string                 //接口描述
	Method string                 //接口请求方式
	Router string                 //接口路径
	Handle func(ctx *gin.Context) //处理方法
	Safe   bool                   //安全接口,不需要验证
}

func StepForward(base string, routers []*WrapReq) {
	if _, ok := ReqManager[base]; !ok {
		ReqManager[base] = make([]*WrapReq, 0)
		ReqManager[base] = append(ReqManager[base], routers...)
	} else {
		ReqManager[base] = append(ReqManager[base], routers...)
	}
	//不需要验证的接口
	filterSafeApi(base, routers)
}

func Loading(engine *gin.Engine) {
	for base, reqs := range ReqManager {
		for _, req := range reqs {
			engine.Handle(req.Method, base+req.Router, req.Handle)
		}
	}
}

func filterSafeApi(base string, routers []*WrapReq) {
	var api string
	for _, router := range routers {
		if router.Safe {
			api = fmt.Sprintf("%s:%s%s", router.Method, base, router.Router)
			if _, ok := SafeApi[api]; !ok {
				SafeApi[api] = router
			}
		}
	}
}

func IsSafeApi(method, path string) bool {
	api := fmt.Sprintf("%s:%s", strings.ToUpper(method), path)
	_, ok := SafeApi[api]
	return ok
}
