package main

import (
	"HelloMyWorld/common/iKafka"
	"HelloMyWorld/config"
	"fmt"
	cluster "github.com/bsm/sarama-cluster"
	"testing"
)

func TestKafka_Consumer(t *testing.T) {
	config.Init()
	iKafka.Init()
	iKafka.Kafka.GroupListenToKafka(config.APPConfig.Kafka.Brokers, "group-log", []string{"log-msg"}, func(cs *cluster.Consumer) {
		for {
			select {
			case msg := <-cs.Messages():
				fmt.Println("Receive message :", string(msg.Value))
			case err := <-cs.Errors():
				fmt.Println("Receive message error :", err.Error())
			}
		}
	})
}
