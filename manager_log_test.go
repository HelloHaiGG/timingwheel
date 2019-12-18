package main

import (
	"HelloMyWorld/common/ikafka"
	"HelloMyWorld/common/ilogger"
	"HelloMyWorld/config"
	"encoding/json"
	"fmt"
	cluster "github.com/bsm/sarama-cluster"
	"testing"
)

func TestManager_Log(t *testing.T) {
	config.Init()
	ikafka.Kafka.GroupListenToKafka(config.APPConfig.Kafka.Brokers, "log-group", []string{config.APPConfig.LTopic.Order}, func(cs *cluster.Consumer) {
		for {
			select {
			case msg := <-cs.Messages():
				m := &ilogger.LogMsg{}
				_ = json.Unmarshal(msg.Value, &m)
				fmt.Printf("%v\n",m)
			case err := <-cs.Errors():
				fmt.Println(err.Error())
			}
		}
	})
}
