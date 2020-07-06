package main

import (
	"HelloMyWorld/common/ielasticsearch"
	"HelloMyWorld/common/iredis"
	"HelloMyWorld/config"
	_ "HelloMyWorld/gateway/app/handles"
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

	config.Init()
	iredis.Init(&iredis.IOptions{
		Host:     config.APPConfig.Redis.Host,
		Port:     config.APPConfig.Redis.Port,
		DB:       config.APPConfig.Redis.DB,
		Password: config.APPConfig.Redis.Password,
	})
	ielasticsearch.Init(&ielasticsearch.IOptions{
		Host:     config.APPConfig.ES.Host,
		Sniff:    config.APPConfig.ES.Sniff,
		User:     config.APPConfig.ES.User,
		Password: config.APPConfig.ES.Password,
		TraceLog: config.APPConfig.ES.TraceLog,
		ErrLog:   config.APPConfig.ES.ErrLog,
	})
	//用户网关服务
	appServer = &http.Server{
		Addr:         "127.0.0.1:59277", //及鲜app
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
