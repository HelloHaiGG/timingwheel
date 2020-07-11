package main

import (
	"HelloMyWorld/common/ikafka"
	"HelloMyWorld/common/ilogger"
	"HelloMyWorld/config"
	"HelloMyWorld/core/listener/kafka"
)

func main() {
	config.Init()
	ikafka.Init(config.APPConfig.Kafka.Brokers)
	ilogger.Init("Kafka Listener")
	 kafka.RegisterListener()
	select {}
}
