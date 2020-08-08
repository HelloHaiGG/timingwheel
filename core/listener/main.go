package main

import (
	"HelloMyWorld/common"
	"HelloMyWorld/common/ikafka"
	"HelloMyWorld/common/ilogger"
	"HelloMyWorld/common/tencentcloud"
	"HelloMyWorld/config"
	"HelloMyWorld/core/listener/kafka"
)

func main() {
	config.Init()
	ikafka.Init(config.APPConfig.Kafka.Brokers)
	ilogger.Init("Kafka Listener")
	tencentcloud.Init(common.SecretId, common.SecretKey)
	kafka.RegisterListener()
}
