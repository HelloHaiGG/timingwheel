package main

import (
	"HelloMyWorld/common/ikafka"
	"HelloMyWorld/common/ilogger"
	"HelloMyWorld/config"
)

func main() {
	config.Init()
	ikafka.Init(config.APPConfig.Kafka.Brokers)
	ilogger.Init("Kafka Listener")
}
