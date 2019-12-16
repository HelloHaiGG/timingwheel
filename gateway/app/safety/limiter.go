package safety

import (
	"context"
	"time"
)

type Limiter struct {
	tokenChan chan struct{} //取令牌通道
	token     struct{}
	sw        bool //开关
	cap       int  //令牌桶容量
	rate      int  //令牌生产速率 单位 s
	cxt       context.Context
	cancel    context.CancelFunc
}

//初始化
func Init(cap, rate int) *Limiter {
	limiter := &Limiter{
		tokenChan: make(chan struct{}, cap),
		token:     struct{}{},
		cap:       cap,
		rate:      rate,
		sw:        false,
	}
	//初始化令牌
	for i := 0; i < cap; i++ {
		limiter.tokenChan <- struct{}{}
	}
	return limiter
}

func (p *Limiter) Start() *Limiter {
	cxt, cancel := context.WithCancel(context.Background())
	p.cancel = cancel
	p.cxt = cxt
	//开始生产
	p.sw = true
	go p.produce(p.cxt)
	return p
}

func (p *Limiter) Stop() {
	if p.sw {
		p.cancel()
	}
}

//生产
func (p *Limiter) produce(cxt context.Context) {
	ticker := time.NewTicker(time.Second)
	var i int
	for p.sw {
		select {
		case <-ticker.C:
			for i = 0; i < p.rate; i++ {
				p.tokenChan <- p.token
			}
		case <-cxt.Done():
			//关闭通道
			close(p.tokenChan)
			ticker.Stop()
		}
	}
}

//取令牌
func (p *Limiter) Pass() bool {
	if !p.sw {
		return true
	}
	if _, ok := <-p.tokenChan; ok {
		return true
	} else {
		//通道关闭
		return false
	}
}
