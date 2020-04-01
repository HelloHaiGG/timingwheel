package serverhandler

import (
	"context"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"go.etcd.io/etcd/clientv3"
	"google.golang.org/grpc/resolver"
)

/**
实现 Resolver 对象
Resolver:服务地址维护对象,通过两个接口方法可以改变维护状态
*/

const (
	EXIST     = 1
	NOT_EXIST = -1
)

type IResolver struct {
	cc   resolver.ClientConn
	c    *clientv3.Client
	addr []resolver.Address
}

func (p *IResolver) ResolveNow(opt resolver.ResolveNowOptions) {

}

func (p *IResolver) Close() {

}

//监控服务地址
func (p *IResolver) watch(prefix string) {
	//prefix 通过服务注册的前缀,可以获取到所有的服务地址  (如: /ISERVER/OrderServer/*** : 127.0.0.1:8888,127.0.0.2:8888)
	state := resolver.State{}

	//etcd中获取 前缀 key 对应的所有值的结果通道
	wChan := p.c.Watch(context.TODO(), prefix, clientv3.WithPrefix())
	for {
		select {
		case resp := <-wChan:
			for _, v := range resp.Events {
				switch v.Type {
				case mvccpb.PUT:
					if p.isExist(string(v.Kv.Value)) == NOT_EXIST {
						p.addr = append(p.addr, resolver.Address{Addr: string(v.Kv.Value)})
					}
				case mvccpb.DELETE:
					if res := p.isExist(string(v.Kv.Value)) ;res != EXIST {
						p.addr = append(p.addr[:res],p.addr[res+1:]...)
					}
				}
			}
		}
		state.Addresses = p.addr
		p.cc.UpdateState(state)
	}
}

func (p *IResolver) isExist(addr string) int32 {
	for k, _ := range p.addr {
		if p.addr[k].Addr == addr {
			return int32(k)
		}
	}
	return NOT_EXIST
}
