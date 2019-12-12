package iKafka

import (
	"HelloMyWorld/config"
	"fmt"
	"github.com/Shopify/sarama"
	cluster "github.com/bsm/sarama-cluster"
	"log"
	"time"
)

var Kafka KafkaCli

type KafkaCli struct {
	//异步 不等待确认
	asyncProducer sarama.AsyncProducer
	//同步 确认之前阻塞
	syncProducer sarama.SyncProducer
	//配置
	cfg *sarama.Config
}

type KafkaMsg struct {
	Topic     string
	Key       string
	Value     string
	Partition int32
	Offset    int64
}

func Init() {
	cfg := sarama.NewConfig()
	//等待所有服务返回相应结果(主从服务器)
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	//随即选择分区
	cfg.Producer.Partitioner = sarama.NewRandomPartitioner
	cfg.Producer.Timeout = time.Duration(config.APPConfig.Kafka.Timeout) * time.Second
	//结果配置 err:默认设为true
	//如果异步生产者 Success 设置为 true,则必须通过channel接收结果,否者会死锁,同步生产者不需要
	cfg.Producer.Return.Successes = true
	cfg.Producer.Return.Errors = true

	//失败后重新尝试
	cfg.Producer.Retry.Max = 3
	//重试间隔
	cfg.Producer.Retry.Backoff = time.Millisecond * 100

	//消费者生产者版本要一致 ！！！
	cfg.Version = sarama.V2_3_0_0
	Kafka = KafkaCli{
		cfg: cfg,
	}
	//初始化异步生产者
	Kafka.newAsyncProducer()
	//初始化同步生存者
	Kafka.newSyncProducer()

	//监听发送结果
	go Kafka.HandResult(func(msg <-chan *sarama.ProducerMessage) {
		suc := &sarama.ProducerMessage{}
		var key []byte
		var value []byte
		for {
			select {
			case suc = <-msg:
				value, _ = suc.Value.Encode()
				key, _ = suc.Key.Encode()
				fmt.Printf("Time：%s Kafka Send Msg Suc!topic=\"%s\";key=\"%s\";value=\"%s\"\n", suc.Timestamp, suc.Topic, key, value)
			}
		}
	})
	//监听发送失败结果
	go Kafka.HandErr(func(errors <-chan *sarama.ProducerError) {
		err := &sarama.ProducerError{}
		var key []byte
		var value []byte
		for {
			select {
			case err = <-errors:
				value, _ = err.Msg.Value.Encode()
				key, _ = err.Msg.Key.Encode()
				fmt.Printf("Time：%s Kafka Send Msg Suc!topic=\"%s\";key=\"%s\";value=\"%s\"!\n", err.Msg.Timestamp, err.Msg.Topic, key, value)
			}
		}
	})
}

//初始化异步生产者
func (p *KafkaCli) newAsyncProducer() {
	producer, err := sarama.NewAsyncProducer(config.APPConfig.Kafka.Brokers, p.cfg)
	if err != nil {
		log.Fatal("异步生产者初始失败")
		return
	}
	p.asyncProducer = producer
}

//初始化同步生产者
func (p *KafkaCli) newSyncProducer() {
	producer, err := sarama.NewSyncProducer(config.APPConfig.Kafka.Brokers, p.cfg)
	if err != nil {
		log.Fatal("同步生产者初始失败")
		return
	}
	p.syncProducer = producer
}

//异步发送Msg
func (p *KafkaCli) ASyncSendMsg(msg *KafkaMsg) {
	//封装成 message
	kafkaMsg := &sarama.ProducerMessage{
		Topic: msg.Topic,
		Key:   sarama.StringEncoder(msg.Key),
		Value: sarama.StringEncoder(msg.Value),
	}

	fmt.Println(p.cfg.Version)

	//发送信息
	p.asyncProducer.Input() <- kafkaMsg
}

//同步发送多个Msg
func (p *KafkaCli) SyncSendMsg(msg []*KafkaMsg) error {
	sMsg := make([]*sarama.ProducerMessage, len(msg))
	for e := range msg {
		sMsg[e].Key = sarama.StringEncoder(msg[e].Key)
		sMsg[e].Value = sarama.StringEncoder(msg[e].Value)
		sMsg[e].Topic = msg[e].Topic
		sMsg[e].Offset = msg[e].Offset
		sMsg[e].Partition = msg[e].Partition
	}
	return p.syncProducer.SendMessages(sMsg)
}

//同步发送单个Msg
func (p *KafkaCli) SyncSendOneMsg(msg *KafkaMsg) (partition int32, offset int64, err error) {
	//封装成 message
	sMsg := &sarama.ProducerMessage{
		Topic:     msg.Topic,
		Key:       sarama.StringEncoder(msg.Key),
		Value:     sarama.StringEncoder(msg.Value),
		Partition: msg.Partition,
		Offset:    msg.Offset,
	}
	return p.syncProducer.SendMessage(sMsg)
}

//消费者 消费组
func (p *KafkaCli) GroupListenToKafka(brokers []string, groupId string, topics []string, handler func(cs *cluster.Consumer)) {
	cfg := cluster.NewConfig()
	cfg.Consumer.Retry.Backoff = 100 * time.Millisecond
	cfg.Consumer.Offsets.Initial = sarama.OffsetNewest
	cfg.Group.Return.Notifications = true
	cfg.Consumer.Return.Errors = true
	cfg.Version = sarama.V2_3_0_0
	c, err := cluster.NewConsumer(brokers, groupId, topics, cfg)
	if err != nil {
		log.Fatal("初始化消费者失败!")
		return
	}
	//
	go func() {
		for {
			select {
			case <-c.Notifications():
			}
		}
	}()
	handler(c)
}

//消费者 非消费组
func (p *KafkaCli) ListenToKafka() {

}

//处理错误
func (p *KafkaCli) HandErr(handler func(<-chan *sarama.ProducerError)) {
	handler(p.asyncProducer.Errors())
}

//处理成功的结果
func (p *KafkaCli) HandResult(handler func(<-chan *sarama.ProducerMessage)) {
	handler(p.asyncProducer.Successes())
}

func (p *KafkaCli) Close() {
	if err := p.asyncProducer.Close(); err != nil {
		p.asyncProducer.AsyncClose()
	}
	_ = p.syncProducer.Close()
}
