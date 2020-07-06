package main

import (
	"HelloMyWorld/common/ikafka"
	"HelloMyWorld/config"
)

func main() {
	config.Init()
	ikafka.Init(config.APPConfig.Kafka.Brokers)

}
