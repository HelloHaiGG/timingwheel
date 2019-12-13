package main

import (
	"HelloMyWorld/common/iKafka"
	"HelloMyWorld/config"
	"fmt"
	cluster "github.com/bsm/sarama-cluster"
	"github.com/sirupsen/logrus"
)

/**
日志整理模块
通过 kafka 将每个微服务产生的log集中存储到mongo/tablestore,方便便利
*/

func main() {
	config.Init()
	iKafka.Kafka.GroupListenToKafka(config.APPConfig.Kafka.Brokers,
		"log-manager",
		[]string{
			config.APPConfig.LTopic.Order,
			config.APPConfig.LTopic.Finance,
			config.APPConfig.LTopic.Gateway,
		},
		func(cs *cluster.Consumer) {
			for {
				select {
				case msg := <-cs.Messages():
					//TODO 将系统日志存储起来
					fmt.Println(msg)
				case err := <-cs.Errors():
					logrus.Errorf("Log Manager Error:%s\n", err.Error())
				}
			}
		})
}
