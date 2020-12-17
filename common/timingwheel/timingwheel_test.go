package timingwheel

import (
	"HelloMyWorld/common/iredis"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/cast"
	"math/rand"
	"os"
	"os/signal"
	"runtime/trace"
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

func TestTimingWheel_AddTask(t *testing.T) {
	f, _ := os.Create("tw.output")
	defer f.Close()
	_ = trace.Start(f)
	defer trace.Stop()
	tw := NewTimingWheel(&TWOptions{ms: time.Second, size: 30, toCache: true, toLoad: true, behavior: &HTaskCache{}})
	_ = tw.Register("task_id", task)
	_ = tw.Register("repeat_task", repeatTask)
	fmt.Println("start...")
	tw.Start()
	fmt.Println("start suc...")
	go func() {
		t := time.NewTicker(time.Millisecond * 300)
		for {
			select {
			case <-t.C:
				t := rand.Int63n(100)
				err := tw.AddTask((&HTask{}).UUid() + time.Now().String(), &Options{
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
