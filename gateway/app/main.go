package main

import (
	"HelloMyWorld/gateway/app/router"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"time"
)


var appServer *http.Server
var eg *errgroup.Group

func main() {
	eg = &errgroup.Group{}

	//用户网关服务
	appServer = &http.Server{
		Addr:         "127.255.10.1:59277", //及鲜app
		Handler:      router.APPRouter(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	eg.Go(func() error {
		return appServer.ListenAndServe()
	})

	if err := eg.Wait(); err != nil {
		log.Fatal("App gateway err:", err)
	}
	//router.UserRouter().Run(":59277")
}
