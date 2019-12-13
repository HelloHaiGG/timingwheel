package main

import (
	"HelloMyWorld/common/ILogger"
	"HelloMyWorld/common/ILogger/collect"
	"HelloMyWorld/config"
	"fmt"
	"testing"
	"time"
)


func TestInitCollector(t *testing.T) {

	//初始化配置
	config.Init()
	ILogger.ToFile = true
	err := collect.InitCollectorAndStart(config.APPConfig.Servers.Order, config.APPConfig.LTopic.Order)
	if err != nil{
		fmt.Println(err)
	}

	time.Sleep(time.Second)
	logger := ILogger.Init(config.APPConfig.Servers.Order)

	logger.Info("This ILogger error")

	<-make(chan bool)
}

func TestInfo(t *testing.T) {
	//初始化日志库
	logger := ILogger.Init("logger server")
	logger.Debug("This ILogger error")
	time.Sleep(time.Second * 2)
}
