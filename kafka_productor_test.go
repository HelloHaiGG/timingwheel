package main

import (
	"HelloMyWorld/common/iKafka"
	"HelloMyWorld/config"
	"testing"
)

func TestKafka_Producer(t *testing.T) {
	config.Init()
	iKafka.Init(config.APPConfig.Kafka.Brokers)
	defer iKafka.Kafka.Close()
	//异步发送数据
	iKafka.Kafka.ASyncSendMsg(&iKafka.KafkaMsg{
		Topic: config.APPConfig.LTopic.Gateway,
		Key:   "",
		Value: "This is Kafka Test!",
	})
	iKafka.Kafka.ASyncSendMsg(&iKafka.KafkaMsg{
		Topic: config.APPConfig.LTopic.Finance,
		Key:   "",
		Value: "This is Kafka Test!",
	})
	iKafka.Kafka.ASyncSendMsg(&iKafka.KafkaMsg{
		Topic: config.APPConfig.LTopic.Order,
		Key:   "",
		Value: "This is Kafka Test!",
	})
	<-make(chan bool)
}
