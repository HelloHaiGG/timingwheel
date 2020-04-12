package timingwheel

import (
	"context"
)

//任务接口-每一个需要被执行的任务,需要实现这个接口
type TimingTask interface {
	//初始化操作
	Init()
	//执行任务
	//Perform(cxt context.Context) error
	Perform() error
	//任务停止时
	OnStop() error
}

//对任务的包装
type WarpTask struct {
	id       string     //任务标识
	isRepeat bool       //是否是需要重复执行的任务
	round    int64      //在第几轮被执行
	slotIdx  int32      //位于第几个时间槽
	task     TimingTask //任务
	ops      *Options   //记录参数

	ctx    context.Context
	cancel context.CancelFunc
}
