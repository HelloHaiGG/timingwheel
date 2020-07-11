package kafka

import (
	"HelloMyWorld/common"
	sms "HelloMyWorld/common/entity/kafka"
	"HelloMyWorld/common/ikafka"
	"HelloMyWorld/common/ilogger"
	"HelloMyWorld/config"
	"fmt"
	cluster "github.com/bsm/sarama-cluster"
	"github.com/gogo/protobuf/proto"
)

//处理通过kafka发送的消息

func RegisterListener() {
	//处理用户注册记录
	ikafka.Kafka.GroupListenToKafka(config.APPConfig.Kafka.Brokers, common.RegisterGroupID, []string{common.RegisterListenerTopic, common.SendSMSTopic}, func(c *cluster.Consumer) {
		var entity proto.Message
		for {
			select {
			case msg := <-c.Messages():
				err := proto.Unmarshal(msg.Value, entity)
				if err != nil {
					ilogger.Ins.Error(err)
				}
				switch entity.(type) {
				case *sms.SendSms:
					fmt.Println(entity)
				default:
					ilogger.Ins.Warn("Unmatch type in kafka group id: ", common.RegisterGroupID)
				}
			case err := <-c.Errors():
				ilogger.Ins.Error(err)
				//TODO
			}
		}
	})
}
