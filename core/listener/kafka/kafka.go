package kafka

import (
	"HelloMyWorld/common"
	sms "HelloMyWorld/common/entity/kafka"
	"HelloMyWorld/common/ikafka"
	"HelloMyWorld/common/ilogger"
	"HelloMyWorld/config"
	"HelloMyWorld/core/listener/handle"
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
				switch string(msg.Key) {
				case common.SendSMSKey:
					entity = &sms.SendSms{}
					str := string(msg.Value)
					_ = proto.UnmarshalText(str, entity)
				default:
				}
			case err := <-c.Errors():
				ilogger.Ins.Error(err)
			}
			handle.Handler(entity)
		}
	})
}
