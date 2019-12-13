package iKafka

import (
	"HelloMyWorld/config"
	"fmt"
	"github.com/Shopify/sarama"
	cluster "github.com/bsm/sarama-cluster"
	"log"
	"sync"
	"time"
)

var Kafka kafkaCli

type kafkaCli struct {
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

func Init(broker []string) {
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
	Kafka = kafkaCli{
		cfg: cfg,
	}
	//初始化异步生产者
	Kafka.newAsyncProducer(broker)
	//初始化同步生存者
	Kafka.newSyncProducer(broker)

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
				fmt.Printf("Time：%s Kafka Send Msg Suc!topic=\"%s\";key=\"%s\";value=\"%s\".\n", suc.Timestamp, suc.Topic, key, value)
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
				fmt.Printf("Time：%s Kafka Send Msg Suc!topic=\"%s\";key=\"%s\";value=\"%s\".\n", err.Msg.Timestamp, err.Msg.Topic, key, value)
			}
		}
	})
}

//初始化异步生产者
func (p *kafkaCli) newAsyncProducer(broker []string) {
	producer, err := sarama.NewAsyncProducer(broker, p.cfg)
	if err != nil {
		log.Fatal("异步生产者初始失败.")
		return
	}
	p.asyncProducer = producer
}

//初始化同步生产者
func (p *kafkaCli) newSyncProducer(broker []string) {
	producer, err := sarama.NewSyncProducer(broker, p.cfg)
	if err != nil {
		log.Fatal("同步生产者初始失败.")
		return
	}
	p.syncProducer = producer
}

//异步发送Msg
func (p *kafkaCli) ASyncSendMsg(msg *KafkaMsg) {
	//封装成 message
	kafkaMsg := &sarama.ProducerMessage{
		Topic:     msg.Topic,
		Key:       sarama.StringEncoder(msg.Key),
		Value:     sarama.StringEncoder(msg.Value),
		Timestamp: time.Now(),
	}
	//发送信息
	p.asyncProducer.Input() <- kafkaMsg
}

//同步发送多个Msg
func (p *kafkaCli) SyncSendMsg(msg []*KafkaMsg) error {
	sMsg := make([]*sarama.ProducerMessage, len(msg))
	for e := range msg {
		sMsg[e].Key = sarama.StringEncoder(msg[e].Key)
		sMsg[e].Value = sarama.StringEncoder(msg[e].Value)
		sMsg[e].Topic = msg[e].Topic
		sMsg[e].Offset = msg[e].Offset
		sMsg[e].Partition = msg[e].Partition
		sMsg[e].Timestamp = time.Now()
	}
	return p.syncProducer.SendMessages(sMsg)
}

//同步发送单个Msg
func (p *kafkaCli) SyncSendOneMsg(msg *KafkaMsg) (partition int32, offset int64, err error) {
	//封装成 message
	sMsg := &sarama.ProducerMessage{
		Topic:     msg.Topic,
		Key:       sarama.StringEncoder(msg.Key),
		Value:     sarama.StringEncoder(msg.Value),
		Partition: msg.Partition,
		Offset:    msg.Offset,
		Timestamp: time.Now(),
	}
	return p.syncProducer.SendMessage(sMsg)
}

//消费者 消费组
func (p *kafkaCli) GroupListenToKafka(brokers []string, groupId string, topics []string, handler func(cs *cluster.Consumer)) {
	cfg := cluster.NewConfig()
	cfg.Consumer.Retry.Backoff = 100 * time.Millisecond
	cfg.Consumer.Offsets.Initial = sarama.OffsetNewest
	cfg.Group.Return.Notifications = true
	cfg.Consumer.Return.Errors = true
	cfg.Version = sarama.V2_3_0_0
	c, err := cluster.NewConsumer(brokers, groupId, topics, cfg)
	if err != nil {
		log.Fatal("初始化消费组失败.")
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
	defer c.Close()
}

//消费者 非消费组
func (p *kafkaCli) ListenToKafka(brokers []string, topic string, handler func(consumer sarama.PartitionConsumer)) {
	cfg := sarama.NewConfig()
	cfg.Consumer.Return.Errors = true
	cfg.Version = sarama.V2_3_0_0

	//创建消费者实例
	c, err := sarama.NewConsumer(brokers, cfg)
	if err != nil {
		log.Fatal("初始化消费者失败.")
	}
	defer c.Close()
	//获取该topic所在的所以分区
	partitions, err := c.Partitions(topic)
	if err != nil {
		log.Fatal("消费者获取分区失败.")
		return
	}
	var pc sarama.PartitionConsumer
	wg := &sync.WaitGroup{}
	//在每个分区拉取消息
	for e := range partitions {
		pc, err = c.ConsumePartition(topic, partitions[e], sarama.OffsetNewest)
		if err != nil {
			continue
		}
		//记录存在消息的分区数量
		wg.Add(1)
		//拉去消息
		go func() {
			handler(pc)
			wg.Done()
			defer pc.AsyncClose()
		}()
	}
	wg.Wait()
}

//处理错误
func (p *kafkaCli) HandErr(handler func(<-chan *sarama.ProducerError)) {
	handler(p.asyncProducer.Errors())
}

//处理成功的结果
func (p *kafkaCli) HandResult(handler func(<-chan *sarama.ProducerMessage)) {
	handler(p.asyncProducer.Successes())
}

func (p *kafkaCli) Close() {
	if err := p.asyncProducer.Close(); err != nil {
		p.asyncProducer.AsyncClose()
	}
	_ = p.syncProducer.Close()
}
