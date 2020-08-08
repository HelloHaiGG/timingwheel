package main

import (
	"HelloMyWorld/common/ikafka"
	"HelloMyWorld/config"
	"fmt"
	"github.com/Shopify/sarama"
	cluster "github.com/bsm/sarama-cluster"
	"testing"
)

func TestKafka_GroupConsumer(t *testing.T) {
	config.Init()
	ikafka.Kafka.GroupListenToKafka(config.APPConfig.Kafka.Brokers, "group-log", []string{"log-msg"}, func(cs *cluster.Consumer) {
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
func TestKafka_Consumer(t *testing.T) {
	config.Init()
	ikafka.Init(config.APPConfig.Kafka.Brokers)
	ikafka.Kafka.ListenToKafka(config.APPConfig.Kafka.Brokers, "Log-Kafka-Topic-GetWay", func(pc sarama.PartitionConsumer) {
		for {
			select {
			case msg := <-pc.Messages():
				fmt.Println("Receive message :", string(msg.Value))
			case err := <-pc.Errors():
				fmt.Println("Receive message error :", err.Error())
			}
		}
	})
}
