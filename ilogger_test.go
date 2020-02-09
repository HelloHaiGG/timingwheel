package main

import (
	"HelloMyWorld/common/ilogger"
	"HelloMyWorld/common/ilogger/collect"
	"HelloMyWorld/config"
	"fmt"
	"testing"
	"time"
)


func TestInitCollector(t *testing.T) {



	//初始化配置
	config.Init()
	ilogger.ToFile = true
	err := collect.InitCollectorAndStart(config.APPConfig.Servers.Order, config.APPConfig.LTopic.Order)
	if err != nil{
		fmt.Println(err)
	}

	time.Sleep(time.Second)
	logger := ilogger.Init(config.APPConfig.Servers.Order)

	logger.Info("This ilogger error")

	<-make(chan bool)
}

func TestInfo(t *testing.T) {
	//初始化日志库
	logger := ilogger.Init("logger server")
	logger.Debug("This ilogger error")
	time.Sleep(time.Second * 2)
}
