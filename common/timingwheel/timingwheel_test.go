package timingwheel

import (
	"fmt"
	"testing"
	"time"
)

type HTask struct {

}

func (p *HTask)Init()  {
	fmt.Println("H Task Init...")
}
func (p *HTask)Perform() error {
	fmt.Println("H Task Perform...")
	return nil
}
func (p *HTask)OnStop() error {
	fmt.Println("H Task Stop...")
	return nil
}

func TestTimingWheel(t *testing.T) {
	tw := NewTimingWheel(time.Second,3)
	tw.Start()
	id,err := tw.AddTask(&HTask{},&Options{
		TimingTime: 10,
		IsRepeat:   true,
	})
	go func() {
		time.Sleep(time.Second * 11)
		if tw.DelTask(id){
			fmt.Println("del success.")
		}else{
			fmt.Println("sss")
		}
	}()
	fmt.Println(id,err)
	for{
		t := time.NewTicker(time.Second)
		select {
		case <-t.C:
			fmt.Println(tw.currentTime)
		}
	}
}
