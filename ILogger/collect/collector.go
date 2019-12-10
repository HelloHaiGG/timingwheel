package collect

import (
	"fmt"
	"github.com/hpcloud/tail"
	"log"
	"os"
)

type Collector struct {
	Path  string
	Topic string
}

func InitCollector(server, topic string) *Collector {
	c := &Collector{}
	if server == "" || topic == "" {
		panic("Collector 'Server' or 'Kafka-Topic' is nil!")
	} else {
		c.Topic = topic
		path, _ := os.Getwd()
		c.Path = path + fmt.Sprintf("\\logs\\%s.txt", server)
	}
	return c
}

func (p *Collector) Start() {
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
				fmt.Println(log.Text)
				//发送到kafka
			default:
			}
		}
	}(t)
}