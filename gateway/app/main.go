package main

import (
	"HelloMyWorld/gateway/app/router"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"time"
)

//TODO APP网关

//1.restful api
//2.限流
//3.接口防刷
//4.权限验证
//5.结果返回
//6.反爬虫
//7.熔断 //go hystrix
//8.ip黑/白名单
//9.灰度发布

var userServer *http.Server
var eg *errgroup.Group

func main() {
	eg = &errgroup.Group{}

	//用户网关服务
	userServer = &http.Server{
		Addr:         ":59277", //及鲜app
		Handler:      router.UserRouter(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	eg.Go(func() error {
		return userServer.ListenAndServe()
	})

	if err := eg.Wait(); err != nil {
		log.Fatal("App gateway err:", err)
	}
	//router.UserRouter().Run(":59277")
}
