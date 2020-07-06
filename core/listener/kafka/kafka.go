package kafka

import (
	"HelloMyWorld/common"
	"HelloMyWorld/common/ikafka"
	"HelloMyWorld/config"
	cluster "github.com/bsm/sarama-cluster"
)

//处理通过kafka发送的消息

func RegisterListener() {
	//处理用户注册记录
	ikafka.Kafka.GroupListenToKafka(config.APPConfig.Kafka.Brokers, "kafka_register_listener", []string{common.RegisterListenerTopic}, func(c *cluster.Consumer) {
		for {
			select {
			case _ = <-c.Messages():
				//TODO
			case _ = <-c.Errors():
				//TODO
			}
		}
	})
}
