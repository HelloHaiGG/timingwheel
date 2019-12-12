package main

import (
	"HelloMyWorld/common/iKafka"
	"HelloMyWorld/config"
	"testing"
)

func TestKafka_Producer(t *testing.T) {
	config.Init()
	iKafka.Init()
	//异步发送数据
	iKafka.Kafka.ASyncSendMsg(&iKafka.KafkaMsg{
		Topic: "log-msg",
		Key:   "",
		Value: "This is Kafka Test!",
	})
	<-make(chan bool)
}
