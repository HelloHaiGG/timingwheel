package timingwheel

import (
	"context"
	"errors"
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
	wheelSize   int32         //槽的数量
	interval    time.Duration //时间轮的跨度
	currentTime int64         //时间轮走针
	slots       []*wrapList   //时间轮的每一个槽
	total       int32         //时间轮的总任务量
	started     bool          //是否已经开启

	//task-map
	taskManager map[string]*WarpTask //管理所有任务
}

//一些设置选项
type Options struct {
	TimingTime time.Duration //执行时间
	IsRepeat   bool          //是否需要重复执行
}

//创建一个时间轮
func NewTimingWheel(ms time.Duration, size int32) *TimingWheel {
	if size <= 0 {
		return nil
	}
	return &TimingWheel{
		tickMs:    ms,
		wheelSize: size,
		interval:  ms * time.Duration(size),
		slots:     make([]*wrapList, size),
		taskManager:make(map[string]*WarpTask),
	}
}

//添加任务
func (p *TimingWheel) AddTask(task TimingTask, opts *Options) (id string, err error) {
	if !p.started {
		err = errors.New("Timing-Wheel Not Working. You Should Start Timing-Wheel First. ")
	}
	if opts == nil {
		err = errors.New("Timing-Wheel Options Not Null! ")
		return
	}
	id = "Task_" + time.Now().String() //TODO 需要替换成UUID
	ctx, cancel := context.WithCancel(context.Background())
	wrapTask := &WarpTask{
		id:       id,
		isRepeat: opts.IsRepeat,
		round:    int32(opts.TimingTime) / p.wheelSize,                             //计算出所在圈
		slotIdx:  int32(p.currentTime + int64(opts.TimingTime)%int64(p.wheelSize)), //计算出所在时间槽
		task:     task,
		ctx:      ctx,
		cancel:   cancel,
	}
	p.addToList(wrapTask)

	//添加到管理器
	p.taskManager[id] = wrapTask
	return
}

//删除任务
func (p *TimingWheel) DelTask(id string) bool {
	if !p.started {
		return false
	}
	if _, ok := p.taskManager[id]; !ok {
		return false
	}
	if p.delTask(id) {
		delete(p.taskManager, id)
		return true
	}
	return false
}

func (p *TimingWheel)Start()  {
	p.started = true
}

func (p *TimingWheel)Stop()  {
	p.started = false
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
