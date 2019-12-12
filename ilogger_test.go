package main

import (
	"HelloMyWorld/common/ILogger"
	"HelloMyWorld/common/ILogger/collect"
	"HelloMyWorld/config"
	"testing"
	"time"
)


func TestInitCollector(t *testing.T) {

	//初始化配置
	config.Init()

	c := collect.InitCollector(config.APPConfig.Servers.Order, config.APPConfig.LTopic.Order)
	c.Start()

	time.Sleep(time.Second)
	logger := ILogger.Init(config.APPConfig.Servers.Order)
	logger.Error("This ILogger error")
	logger.Info("This ILogger error")
	logger.Info("This ILogger error")
	logger.Info("This ILogger error")
	logger.Info("This ILogger error")
	logger.Info("This ILogger error")
	logger.Info("This ILogger error")
	logger.Info("This ILogger error")
	logger.Info("This ILogger error")
	logger.Info("This ILogger error")
	logger.Info("This ILogger error")
	logger.Info("This ILogger error")
	logger.Info("This ILogger error")
	logger.Info("This ILogger error")
	logger.Info("This ILogger error")
	logger.Info("This ILogger error")
	logger.Info("This ILogger error")
	logger.Info("This ILogger error")
	logger.Info("This ILogger error")
	logger.Info("This ILogger error")
	logger.Info("This ILogger error")

	<-make(chan bool)
}

func TestInfo(t *testing.T) {
	//初始化日志库
	logger := ILogger.Init("logger server")
	logger.Debug("This ILogger error")
	logger.Error("This ILogger error")
	logger.Info("This ILogger error")
	time.Sleep(time.Second * 2)
}
