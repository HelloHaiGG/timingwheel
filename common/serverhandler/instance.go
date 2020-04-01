package serverhandler

import (
	"HelloMyWorld/config"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"log"
	"time"
)

func GetServerConn(server string) *grpc.ClientConn  {
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithBalancerName(roundrobin.Name),
		//grpc.WithBlock(),
	}
	cxt,_ := context.WithTimeout(context.Background(),time.Duration(config.APPConfig.Grpc.CallTimeOut)*time.Second)
	conn ,err := grpc.DialContext(cxt,DNSName+":///"+server,opts...)
	if err != nil{
		log.Fatalf("Grpc Dial %s Error:%v",server,err)
	}
	return conn
}
