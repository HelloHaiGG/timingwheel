package collect

import (
	"HelloMyWorld/common/ilogger"
	"HelloMyWorld/common/ikafka"
	"HelloMyWorld/config"
	"errors"
	"fmt"
	"github.com/hpcloud/tail"
	"log"
	"os"
)

type Collector struct {
	Path  string
	Topic string
}

func init() {
	config.Init()
	//初始化kafka
	ikafka.Init(config.APPConfig.Kafka.Brokers)
}

func InitCollectorAndStart(server, topic string) error {

	//本地测试环境默认不开启日志收集功能
	if !ilogger.ToFile {
		return errors.New("Start Collector failed. ilogger Config of 'ToFile' is 'false' ")
	}
	c := &Collector{}
	if server == "" || topic == "" {
		panic("Collector 'Server' or 'Kafka-Topic' is nil!")
	} else {
		c.Topic = topic
		path, _ := os.Getwd()
		c.Path = path + fmt.Sprintf("/logs/%s.txt", server)
	}
	c.start()
	return nil
}

func (p *Collector) start() {
	t, err := tail.TailFile(p.Path, tail.Config{
		Location: &tail.SeekInfo{
			Offset: 0,
			Whence: os.SEEK_END,
		},
		Follow:    true,
		ReOpen:    true,
		Poll:      true,
		MustExist: false,
	})
	if err != nil {
		log.Printf("load collecter failed from path: %s || ERR: %s", p.Path, err)
		return
	}
	go func(tail *tail.Tail) {
		for {
			select {
			case log := <-t.Lines:
				//发送到kafka
				ikafka.Kafka.ASyncSendMsg(&ikafka.KafkaMsg{
					Topic: p.Topic,
					Value: log.Text,
				})
			}
		}
	}(t)
}
