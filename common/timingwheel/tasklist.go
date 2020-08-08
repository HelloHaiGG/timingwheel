package timingwheel

import (
	"sort"
	"sync"
)

type taskList []*WrapTask

type wrapList struct {
	total      int32    //列表中的任务量
	sync.Mutex          //对任务列表进行保护
	list       taskList // 列表
}

//向队列中添加任务
func (p *wrapList) add(task *WrapTask) {
	if p.list == nil {
		p.list = make([]*WrapTask, 0)
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
	var task *WrapTask
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
func (p *wrapList) get() []*WrapTask {
	if p == nil || p.list == nil {
		return nil
	}
	var tasks []*WrapTask
	tasks = p.CatOutTask()
	for i, _ := range p.list {
		p.list[i].round--
	}
	p.total = int32(len(p.list))
	for _, task := range tasks {
		TWCline.taskManager.Delete(task.id)
		if task.isRepeat {
			//重复加入
			if err := TWCline.AddTask(task.id, task.ops); err != nil {
				//TODO 重复加入失败
			}
		}
	}
	return tasks
}

func (p *wrapList) CatOutTask() []*WrapTask {
	if p == nil || p.list == nil || len(p.list) <= 0 {
		return nil
	}
	//全部未到执行时间
	if p.list[0].round > 0 {
		return nil
	}
	index := -1
	for i, task := range p.list {
		if task.round == 0 {
			index = i
		} else {
			break
		}
	}
	if index > -1 {
		sub := p.list[:index+1]
		p.list = p.list[index+1:]
		return sub
	}
	return nil
}

//任务列表排序
type sortTask []*WrapTask

func (p sortTask) Len() int {
	return len(p)
}

func (p sortTask) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p sortTask) Less(i, j int) bool {
	return p[i].round < p[j].round
}
