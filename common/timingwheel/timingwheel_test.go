package timingwheel

import (
	"errors"
	"fmt"
	"math/rand"
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
	return errors.New("This Error. ")
}
func (p *HTask) OnStop() error {
	fmt.Println("H Task Stop...")
	return nil
}

func TestTimingWheel(t *testing.T) {
	tw := NewTimingWheel(time.Second, 10)
	tw.Start()
	for i := 0; i < 100; i++ {
		t := rand.Int63n(1000)
		id, err := tw.AddTask(&HTask{}, &Options{
			TimingTime:    t,
			IsRepeat:      false,
			NeedHandleErr: true,
		})
		fmt.Println(i, id, err)
	}
	for {
		t := time.NewTicker(time.Second)
		select {
		case <-t.C:
			fmt.Println(tw.currentTime)
		case err := <-tw.handleErr:
			fmt.Println(err)
		}
	}
}

func TestTimingWheel_AddTask(t *testing.T) {
	tw := NewTimingWheel(time.Second, 10)
	tw.Start()
	go func() {
		for {
			t := time.NewTicker(time.Second * 30)
			select {
			case <-t.C:
				t := rand.Int63n(10)
				id, err := tw.AddTask(&HTask{}, &Options{
					TimingTime:    t,
					IsRepeat:      false,
					NeedHandleErr: true,
				})
				fmt.Println("Auto Add:", t, id, err)
			}
		}
	}()
	id, err := tw.AddTask(&HTask{}, &Options{
		TimingTime: 8,
		IsRepeat:   true,
	})
	fmt.Println(id, err)
	select {
	}
}
