package timingwheel

import (
	"context"
	"fmt"
	"testing"
	"time"
)

type HTask struct {

}

func (p *HTask)Init()  {
	fmt.Println("H Task Init...")
}
func (p *HTask)Perform(cxt context.Context) error {
	fmt.Println("H Task Perform...")
	return nil
}
func (p *HTask)OnStop() error {
	fmt.Println("H Task Stop...")
	return nil
}

func TestTimingWheel(t *testing.T) {
	tw := NewTimingWheel(time.Second,24)
	tw.Start()
	id,err := tw.AddTask(&HTask{},&Options{
		TimingTime: 10,
		IsRepeat:   false,
	})
	if tw.DelTask(id){
		fmt.Println("del success.")
	}
	fmt.Println(id,err)
}
