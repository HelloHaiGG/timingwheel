package timingwheel

import (
	"fmt"
	"testing"
	"time"
)

type HTask struct {
}

func (p *HTask) Init() {
	fmt.Println("H Task Init...")
}
func (p *HTask) Perform() error {
	fmt.Println("H Task Perform...")
	return nil
}
func (p *HTask) OnStop() error {
	fmt.Println("H Task Stop...")
	return nil
}

func TestTimingWheel(t *testing.T) {
	tw := NewTimingWheel(time.Second, 10)
	tw.Start()
	for i := 0; i < 100; i++ {
		id, err := tw.AddTask(&HTask{}, &Options{
			TimingTime: 1,
			IsRepeat:   false,
		})
		fmt.Println(i,id, err)
	}
	for {
		t := time.NewTicker(time.Second)
		select {
		case <-t.C:
			fmt.Println(tw.currentTime)
		}
	}
}
