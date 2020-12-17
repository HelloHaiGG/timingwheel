package timingwheel

import (
	"context"
)

//任务接口-每一个需要被执行的任务,需要实现这个接口
type TimingTask interface {
	//uid
	UUid() string
}

//一些设置选项
type Options struct {
	TimingTime    int64 //执行时间
	IsRepeat      bool  //是否需要重复执行
	NeedHandleErr bool  //是否需要处理错误信息
}

//对任务的包装
type WrapTask struct {
	id       string   //任务标识
	isRepeat bool     //是否是需要重复执行的任务
	round    int64    //在第几轮被执行
	slotIdx  int32    //位于第几个时间槽
	ops      *Options //参数

	ctx    context.Context
	cancel context.CancelFunc

	createAt int64 //任务创建时间
}
