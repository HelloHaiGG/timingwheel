package timingwheel

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"
)

// 普通时间轮
// 在普通时间轮中,如果任务的触发时间大于了时间轮的跨度,可以通过记录触发圈数来判定是否需要执行任务
// 这种情况下,需要将处于同一个时间槽中的任务遍历一边,判断是否需要执行,不需要执行的任务需要将圈数进行减一操作

// 层级时间轮
// 分层时间轮可以避免这种操作,层级时间轮中的每一个任务都只需要执行的任务,
// 当时针指向一个卡槽之后,只要判断改任务是需要立即执行,还是继续向下层时间轮分发任务
// 对于任务的添加也是一样,需要判断任务的触发时间是否大于当前时间轮的跨度,如果大于,创建(在不存在的情况下)一个上层时间轮,
// 跨度等于 (wheelSize * tickMs) * wheelSize

// 此处使用普通时间轮实现

var one sync.Once

type TimingWheel struct {
	tickMs      time.Duration //粒度,每一个槽的时间跨度 秒
	wheelSize   int64         //槽的数量
	interval    time.Duration //时间轮的跨度
	currentTime int64         //时间轮走针
	slots       []*wrapList   //时间轮的每一个槽
	total       int32         //时间轮的总任务量
	started     bool          //是否已经开启

	//任务注册器,所有的执行任务都必须加载到注册器中,否则无法执行
	registry map[string]func() error

	cxt    context.Context
	cancel context.CancelFunc

	timer *time.Ticker //驱动走针
	//task-map
	taskManager sync.Map //管理所有任务
	//err-handle
	handleErr chan WrapError

	toCache bool //是否需要此久化
	toLoad  bool //是否需要加载持久化的数据

	behavior CacheBehavior

	loadingSuc bool
	locker     *sync.Mutex
	loading    *sync.Cond
}

type TWOptions struct {
	ms       time.Duration
	size     int64
	toCache  bool
	toLoad   bool
	behavior CacheBehavior
}

func (p *TWOptions) init() {
	if p.ms == 0 {
		p.ms = 1
	}
	if p.size == 0 {
		p.size = 100
	}
}

type CacheBehavior interface {
	CacheKey() string
	OnStop(tasks []*CacheTask, pointer int64) error
	OnStart() []*CacheTask
	ReloadPointer() int64
}

//任务持久化结构
type CacheTask struct {
	Id             string `json:"id"`
	SaveTimeAt     int64  `json:"save_time_at"` //记录持久化的时间
	CreateAt       int64  `json:"create_at"`
	IsRepeat       bool   `json:"is_repeat"`        //是否是循环执行的任务
	NeedHandlerErr bool   `json:"need_handler_err"` //该任务是否需要抛出异常
	TimingTime     int64  `json:"timing_time"`      //任务延迟时间
}

//异常封装
type WrapError struct {
	TaskId string
	Err    error
}

var TWCline *TimingWheel

//创建一个时间轮
func NewTimingWheel(opt *TWOptions) *TimingWheel {
	opt.init()
	one.Do(func() {
		ctx, cancel := context.WithCancel(context.Background())
		if TWCline == nil {
			locker := &sync.Mutex{}
			TWCline = &TimingWheel{
				tickMs:      opt.ms,
				wheelSize:   opt.size,
				currentTime: 0,
				interval:    opt.ms * time.Duration(opt.size),
				slots:       make([]*wrapList, opt.size),
				timer:       time.NewTicker(opt.ms),
				cxt:         ctx,
				cancel:      cancel,
				handleErr:   make(chan WrapError),
				registry:    make(map[string]func() error),
				toCache:     opt.toCache,
				toLoad:      opt.toLoad,
				behavior:    opt.behavior,
				locker:      locker,
				loading:     sync.NewCond(locker),
			}
		}
	})
	return TWCline
}

//添加任务
func (p *TimingWheel) AddTask(task string, opts *Options) (err error) {
	if !p.started {
		err = errors.New("Timing-Wheel Not Working. You Should Start Timing-Wheel First. ")
	}
	if opts == nil {
		err = errors.New("Timing-Wheel Options Is Null! ")
		return
	}
	if !p.loadingSuc {
		log.Println("witing loading cache. ")
		p.loading.L.Lock()
		defer p.loading.L.Unlock()
		p.loading.Wait()
		log.Println("get loading cache over. ")
	}
	//判断任务是否存在
	if _, ok := p.taskManager.Load(task); ok {
		log.Println("Add a task with existing. ")
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	wrapTask := &WrapTask{
		id:       task,
		isRepeat: opts.IsRepeat,
		round:    opts.TimingTime / p.wheelSize,                          //计算出所在圈
		slotIdx:  int32((p.currentTime + opts.TimingTime) % p.wheelSize), //计算出所在时间槽
		ctx:      ctx,
		cancel:   cancel,
		ops:      opts,
		createAt: time.Now().Unix(),
	}
	p.addToList(wrapTask)
	//添加到管理器
	p.taskManager.Store(task, wrapTask)
	return
}

var TaskAlreadyExists = errors.New("Task already exists. ")
var TaskDoesNotExist = errors.New("Task not exist. ")

//注册任务到时间轮
func (p *TimingWheel) Register(uuid string, task func() error) error {
	if _, ok := p.registry[uuid]; ok {
		return TaskAlreadyExists
	}
	p.registry[uuid] = task
	return nil
}

//加载任务
func (p *TimingWheel) perform(uuid string) error {
	if task, ok := p.registry[uuid]; !ok {
		return TaskDoesNotExist
	} else {
		return task()
	}
}

//删除任务
func (p *TimingWheel) DelTask(id string) bool {
	if !p.started {
		return false
	}
	if task, ok := p.taskManager.Load(id); !ok {
		return false
	} else {
		task.(*WrapTask).cancel()
	}
	if p.delTask(id) {
		p.taskManager.Delete(id)
		return true
	}
	return false
}

func (p *TimingWheel) Start() {
	p.started = true
	if p.toLoad {
		p.loadCache()
		log.Println("load cache task over.")
	}
	//时间轮指针驱动
	go p.ticker()
}

func (p *TimingWheel) Stop() {
	p.started = false
	p.cancel()
	p.timer.Stop()
	if p.toCache {
		p.cache()
	}
}

func (p *TimingWheel) ticker() {
	var wg sync.WaitGroup
	var err error
	for {
		select {
		case <-p.timer.C:
			p.currentTime = (p.currentTime) % p.wheelSize
			//获取到走针指向的槽的到期任务
			tasks := p.slots[p.currentTime].get()
			p.currentTime++
			//执行任务
			for _, v := range tasks {
				wg.Add(1)
				go func(ctx context.Context) { //TODO 需要使用协程池
					select {
					case <-ctx.Done():
						wg.Done()
						return
					default:
						if err = p.perform(v.id); err != nil && v.ops.NeedHandleErr {
							p.handleErr <- WrapError{v.id, err}
						}
						wg.Done()
					}
				}(v.ctx)
			}
		case <-p.cxt.Done():
			//等待所有任务执行结束
			wg.Wait()
			return
		}
	}
}

func (p *TimingWheel) cache() {
	log.Println("Cache task to redis")
	tasks := make([]*CacheTask, 0)
	p.taskManager.Range(func(key, value interface{}) bool {
		tasks = append(tasks, &CacheTask{
			Id:             value.(*WrapTask).id,
			SaveTimeAt:     time.Now().Unix(),
			CreateAt:       value.(*WrapTask).createAt,
			IsRepeat:       value.(*WrapTask).isRepeat,
			NeedHandlerErr: value.(*WrapTask).ops.NeedHandleErr,
			TimingTime:     value.(*WrapTask).ops.TimingTime,
		})
		return true
	})
	p.loadingSuc = false
	err := p.behavior.OnStop(tasks, p.currentTime)
	if err != nil {
		log.Println("Timing-Wheel cache to task err: ", err)
	}
	log.Println("Cache task to redis suc.")
}

func (p *TimingWheel) loadCache() {
	log.Println("load cache task")
	tasks := p.behavior.OnStart()
	p.currentTime = p.behavior.ReloadPointer()
	now := time.Now().Unix()
	go func() {
		for _, task := range tasks {
			//已经过期任务立即执行
			if task.SaveTimeAt+task.TimingTime <= now && !task.IsRepeat {
				if do, ok := p.registry[task.Id]; !ok {
					log.Println("Timing wheel load task", task.Id, " fail. not exits. ")
				} else {
					if err := do(); err != nil && task.NeedHandlerErr {
						p.handleErr <- WrapError{Err: err, TaskId: task.Id}
					}
				}
				continue
			}
			//已经过期且任务为循环执行的任务
			if task.SaveTimeAt+task.TimingTime <= now && task.IsRepeat {
				times := (now - task.SaveTimeAt) / task.TimingTime
				if do, ok := p.registry[task.Id]; !ok {
					log.Println("Timing wheel load task", task.Id, " fail. not exits. ")
				} else {
					for i := 0; i < int(times); i++ {
						if err := do(); err != nil && task.NeedHandlerErr {
							p.handleErr <- WrapError{Err: err, TaskId: task.Id}
						}
					}
				}
			}
			// 未过期的任务添加到时间轮
			ctx, cancel := context.WithCancel(context.Background())
			t := &WrapTask{
				id:       task.Id,
				isRepeat: task.IsRepeat,
				round:    task.TimingTime / p.wheelSize,                                                //计算出所在圈
				slotIdx:  int32((p.currentTime + (now-task.SaveTimeAt)%task.TimingTime) % p.wheelSize), //计算出所在时间槽
				ctx:      ctx,
				cancel:   cancel,
				ops: &Options{
					TimingTime:    task.TimingTime,
					IsRepeat:      task.IsRepeat,
					NeedHandleErr: task.NeedHandlerErr,
				},
				createAt: time.Now().Unix(),
			}
			p.addToList(t)
			p.taskManager.Store(task.Id, t)
		}
		defer func() {
			p.loading.Signal()
			p.loadingSuc = true
		}()
	}()
}

//将任务添加到任务列表
//先将任务添加到任务列表(取锁)
func (p *TimingWheel) addToList(task *WrapTask) {
	if p.slots[task.slotIdx] == nil {
		p.slots[task.slotIdx] = new(wrapList)
	}
	p.slots[task.slotIdx].add(task)
}

func (p *TimingWheel) delTask(id string) bool {
	task, ok := p.taskManager.Load(id)
	if !ok {
		return false
	}
	if p.slots[task.(*WrapTask).slotIdx] == nil {
		return false
	}
	return p.slots[task.(*WrapTask).slotIdx].del(id)
}
