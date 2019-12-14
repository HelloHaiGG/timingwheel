package main

import (
	"HelloMyWorld/gateway/app/router"
	"log"
	"net/http"
	"sync"
)

//TODO APP网关

//1.restful api
//2.限流
//3.接口防刷
//4.权限验证
//5.结果返回
//6.反爬虫
//7.熔断
//8.ip黑/白名单
//9.灰度发布

var userServer *http.Server
var wg *sync.WaitGroup

func main() {
	wg = &sync.WaitGroup{}

	//用户网关服务
	userServer = &http.Server{
		Addr:         ":59277", //及鲜app
		Handler:      router.UserRouter(),
		ReadTimeout:  10,
		WriteTimeout: 10,
	}
	wg.Add(1)
	go func() {
		err := userServer.ListenAndServe()
		if err != nil {
			log.Fatal("user router start err:", err)
			wg.Done()
		}
	}()

	wg.Wait()
}
