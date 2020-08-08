package timingwheel

import (
	"HelloMyWorld/common/iredis"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/cast"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"
)

func init() {
	iredis.Init(&iredis.IOptions{
		Host:     "182.92.239.63",
		Port:     6379,
		DB:       0,
		Password: "root",
	})
}

type HTask struct {
}

func (p *HTask) UUid() string {
	return "task_id"
}

type RTask struct {
}

func (p *RTask) UUid() string {
	return "repeat_task"
}

type HTaskCache struct {
}

func (p *HTaskCache) CacheKey() string {
	return "sss"
}

func (p *HTaskCache) OnStop(tasks []*CacheTask, pointer int64) error {
	str, _ := jsoniter.MarshalToString(tasks)
	if _, err := iredis.RedisCli.Set(fmt.Sprintf("%s:TEMING_WHEEL_TASK", p.CacheKey()), str, -1).Result(); err != nil {
		return err
	}
	_, err := iredis.RedisCli.Set(fmt.Sprintf("%s:TEMING_WHEEL_POINTER", p.CacheKey()), pointer, -1).Result()
	return err
}

func (p *HTaskCache) OnStart() []*CacheTask {
	tasks := make([]*CacheTask, 0)
	str, _ := iredis.RedisCli.Get(fmt.Sprintf("%s:TEMING_WHEEL_TASK", p.CacheKey())).Result()
	_ = jsoniter.UnmarshalFromString(str, &tasks)
	return tasks
}

func (p *HTaskCache) ReloadPointer() int64 {
	str, _ := iredis.RedisCli.Get(fmt.Sprintf("%s:TEMING_WHEEL_POINTER", p.CacheKey())).Result()
	return cast.ToInt64(str)
}

func task() error {
	fmt.Println("this is a task.")
	return nil
}

func repeatTask() error {
	fmt.Println("this is a repeat task.")
	return nil
}

func TestTimingWheel(t *testing.T) {
	tw := NewTimingWheel(&TWOptions{ms: 1, size: 10})
	tw.Start()
	for i := 0; i < 100; i++ {
		t := rand.Int63n(1000)
		err := tw.AddTask((&HTask{}).UUid(), &Options{
			TimingTime:    t,
			IsRepeat:      false,
			NeedHandleErr: true,
		})
		fmt.Println(i, err)
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
	tw := NewTimingWheel(&TWOptions{ms: time.Second, size: 30, toCache: true, toLoad: true, behavior: &HTaskCache{}})
	_ = tw.Register("task_id", task)
	_ = tw.Register("repeat_task", repeatTask)
	tw.Start()
	go func() {
		t := time.NewTicker(time.Second * 30)
		for {
			select {
			case <-t.C:
				t := rand.Int63n(10)
				err := tw.AddTask((&HTask{}).UUid(), &Options{
					TimingTime:    t,
					IsRepeat:      false,
					NeedHandleErr: true,
				})
				fmt.Println("Auto Add:", t, err)
			}
		}
	}()
	_ = tw.AddTask((&RTask{}).UUid(), &Options{
		TimingTime: 8,
		IsRepeat:   true,
	})
	stopChan := make(chan os.Signal)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)
	<-stopChan // wait for SIGINT or SIGTERM
	tw.Stop()
}
