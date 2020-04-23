package timingwheel

import (
	"context"
	"errors"
	"log"
	"time"
)

// 普通时间轮
// 在普通时间轮中,如果任务的触发时间大于了时间轮的跨度,可以通过记录触发圈数来判定是否需要执行任务
// 这种情况下,需要将处于同一个时间槽中的任务遍历一边,判断是否需要执行,不需要执行的任务需要多圈数进行减一操作

// 层级时间轮
// 分层时间轮可以避免这种操作,层级时间轮中的没一个任务都只需要执行的任务,
// 当时针指向一个卡槽之后,只要判断改任务是需要立即执行,还是继续向下层时间轮分发任务
// 对于任务的添加也是一样,需要判断任务的触发时间是否大于当前时间轮的跨度,如果大于,创建(在不存在的情况下)一个上层时间轮,
// 跨度等于 (wheelSize * tickMs) * wheelSize

// 此处使用普通时间轮实现

type TimingWheel struct {
	tickMs      time.Duration //粒度,每一个槽的时间跨度 秒
	wheelSize   int64         //槽的数量
	interval    time.Duration //时间轮的跨度
	currentTime int64         //时间轮走针
	slots       []*wrapList   //时间轮的每一个槽
	total       int32         //时间轮的总任务量
	started     bool          //是否已经开启

	cxt    context.Context
	cancel context.CancelFunc

	timer *time.Ticker //驱动走针
	//task-map
	taskManager map[string]*WarpTask //管理所有任务
}

//一些设置选项
type Options struct {
	TimingTime time.Duration //执行时间
	TaskId     string        //任务ID
	IsRepeat   bool          //是否需要重复执行
}

var TWCline *TimingWheel

//创建一个时间轮
func NewTimingWheel(ms time.Duration, size int64) *TimingWheel {
	if size <= 0 {
		return nil
	}
	ctx, cancel := context.WithCancel(context.Background())
	if TWCline == nil {
		TWCline = &TimingWheel{
			tickMs:      ms,
			wheelSize:   size,
			currentTime: 0,
			interval:    ms * time.Duration(size),
			slots:       make([]*wrapList, size),
			taskManager: make(map[string]*WarpTask),
			timer:       time.NewTicker(ms),
			cxt:         ctx,
			cancel:      cancel,
		}
	}
	return TWCline
}

//添加任务
func (p *TimingWheel) AddTask(task TimingTask, opts *Options) (id string,err error) {
	if !p.started {
		err = errors.New("Timing-Wheel Not Working. You Should Start Timing-Wheel First. ")
	}
	if opts == nil {
		err = errors.New("Timing-Wheel Options Not Null! ")
		return
	}
	if opts.TaskId == "" {
		opts.TaskId = "Task_" + time.Now().String() //TODO 需要替换成UUID
	}
	id = opts.TaskId
	ctx, cancel := context.WithCancel(context.Background())
	wrapTask := &WarpTask{
		id:       opts.TaskId,
		isRepeat: opts.IsRepeat,
		round:    int64(opts.TimingTime) / p.wheelSize,                          //计算出所在圈
		slotIdx:  int32((p.currentTime + int64(opts.TimingTime)) % p.wheelSize), //计算出所在时间槽
		task:     task,
		ctx:      ctx,
		cancel:   cancel,
		ops:      opts,
	}
	p.addToList(wrapTask)

	//添加到管理器
	p.taskManager[opts.TaskId] = wrapTask
	return
}

//删除任务
func (p *TimingWheel) DelTask(id string) bool {
	if !p.started {
		return false
	}
	if task, ok := p.taskManager[id]; !ok {
		return false
	} else {
		_ = task.task.OnStop()
		task.cancel()
	}
	if p.delTask(id) {
		delete(p.taskManager, id)
		return true
	}
	return false
}

func (p *TimingWheel) Start() {
	p.started = true
	//时间轮指针驱动
	go p.ticker()
}

func (p *TimingWheel) Stop() {
	p.started = false
	p.cancel()
	p.timer.Stop()
}

func (p *TimingWheel) ticker() {
	for {
		select {
		case <-p.timer.C:
			p.currentTime = (p.currentTime + 1) % p.wheelSize
			//获取到指针指向的槽的到期任务
			tasks := p.slots[p.currentTime].get()
			//执行任务
			for _, v := range tasks {
				go func(ctx context.Context) { //TODO 需要使用协程池
					select {
					case <-ctx.Done():
						return
					default:
						if err := v.task.Perform(); err != nil {
							log.Printf("Task perform err: %v", err)
						}
					}
				}(v.ctx)
			}
		case <-p.cxt.Done():
			return
		}
	}
}

//将任务添加到任务列表
//先将任务添加到任务列表(取锁)
func (p *TimingWheel) addToList(task *WarpTask) {
	if p.slots[task.slotIdx] == nil {
		p.slots[task.slotIdx] = new(wrapList)
	}
	p.slots[task.slotIdx].add(task)
}

func (p *TimingWheel) delTask(id string) bool {
	task := p.taskManager[id]
	if p.slots[task.slotIdx] == nil {
		return false
	}
	return p.slots[task.slotIdx].del(id)
}
