package main

import (
	"HelloMyWorld/common/ikafka"
	"HelloMyWorld/config"
	"testing"
)

func TestKafka_Producer(t *testing.T) {
	config.Init()
	ikafka.Init(config.APPConfig.Kafka.Brokers)
	defer ikafka.Kafka.Close()
	//异步发送数据
	ikafka.Kafka.ASyncSendMsg(&ikafka.KafkaMsg{
		Topic: config.APPConfig.LTopic.Gateway,
		Key:   "",
		Value: "This is Kafka Test!",
	})
	ikafka.Kafka.ASyncSendMsg(&ikafka.KafkaMsg{
		Topic: config.APPConfig.LTopic.Finance,
		Key:   "",
		Value: "This is Kafka Test!",
	})
	ikafka.Kafka.ASyncSendMsg(&ikafka.KafkaMsg{
		Topic: config.APPConfig.LTopic.Order,
		Key:   "",
		Value: "This is Kafka Test!",
	})
	<-make(chan bool)
}
