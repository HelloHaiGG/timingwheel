package timingwheel

import (
	"sort"
	"sync"
)

type taskList []*WarpTask

type wrapList struct {
	total      int32    //列表中的任务量
	sync.Mutex          //对任务列表进行保护
	list       taskList // 列表
}

//向队列中添加任务
func (p *wrapList) add(task *WarpTask) {
	if p.list == nil {
		p.list = make([]*WarpTask, 0)
	}
	p.Lock()
	defer p.Unlock()
	p.list = append(p.list, task)
	p.total++
	//排序,保证最快被执行的任务位于队列前面
	sort.Sort(sortTask(p.list))
}

//删除队列中的任务
func (p *wrapList) del(id string) bool {
	var idx int
	var task *WarpTask
	for idx, task = range p.list {
		if task.id == id {
			break
		}
	}
	if task == nil {
		return false
	}
	p.Lock()
	defer p.Unlock()
	if task.round == 0 {
		//取消任务
		task.cancel()
	}
	p.list = append(p.list[:idx], p.list[idx+1:]...)
	return true
}

//在任务队列中取出到期的任务
func (p *wrapList) get() []*WarpTask {
	if p == nil || p.list == nil {
		return nil
	}
	var tasks []*WarpTask
	var flag bool //判断是否已经取出任务
	for i, task := range p.list {
		//因为任务放入队列后,对任务进行了排序,所以能保证队列前面的任务都是需要立即执行的
		if task.round != 0 || (task.round == 0 && i == len(p.list)-1) {
			if !flag {
				//取出要执行的任务
				tasks = append(tasks, p.list[:i+1]...)
				//剩余的继续留在槽内
				p.list = p.list[i+1:]
				delete(TWCline.taskManager, task.id)
				//判断是否是需要重复执行的任务
				if task.isRepeat {
					//重复加入
					if _, err := TWCline.AddTask(task.task, task.ops); err != nil {
						//TODO 重复加入失败
					}
				}
				//标记已经取出任务
				flag = true
			} else {
				//对剩余任务 round 减一
				p.list[i].round--
			}
		}
	}
	p.total = int32(len(p.list))
	return tasks
}

//任务列表排序
type sortTask []*WarpTask

func (p sortTask) Len() int {
	return len(p)
}

func (p sortTask) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p sortTask) Less(i, j int) bool {
	return p[i].round < p[j].round
}
