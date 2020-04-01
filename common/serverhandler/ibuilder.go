package serverhandler

import (
	"go.etcd.io/etcd/clientv3"
	"google.golang.org/grpc/resolver"
	"log"
)

/**
基于 GRPC 提供的 resolver 接口实现服务发现与注册,并提供了负载均衡策略
IBuilder 创建一个解析器
*/

func init() {
	//注册解析器
	resolver.Register(NewIBuilder())
}

const DNSName  = ""

// IBuilder实现 实现Builder接口
type IBuilder struct {
	DNS string
}

func NewIBuilder() *IBuilder {
	return &IBuilder{
		DNS:     DNSName,
	}
}

// GRPC 将所有配置的 Resolver 维护到一个ResolverMap  --> map[string]Resolver
// GRPC 通过 Scheme() 获取 Resolver 的 key
func (p *IBuilder) Scheme() string {
	return p.DNS
}

func (p *IBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	var etcdClient *clientv3.Client
	var err error
	if etcdClient,err = clientv3.New(clientv3.Config{
		Endpoints:            nil,
		AutoSyncInterval:     0,
		DialTimeout:          0,
		DialKeepAliveTime:    0,
		DialKeepAliveTimeout: 0,
		MaxCallSendMsgSize:   0,
		MaxCallRecvMsgSize:   0,
		TLS:                  nil,
		Username:             "",
		Password:             "",
		RejectOldCluster:     false,
		DialOptions:          nil,
		LogConfig:            nil,
		Context:              nil,
		PermitWithoutStream:  false,
	});err != nil{
		log.Fatalf("Register Resolver Builder Fail:%v",err)
		return nil,err
	}

	//创建 IResolver
	iResolver := &IResolver{
		cc: cc,
		c:  etcdClient,
	}

	go iResolver.watch("/"+target.Scheme + "/" +target.Endpoint)
}






