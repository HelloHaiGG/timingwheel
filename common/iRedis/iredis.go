package iRedis

import (
	"fmt"
	"github.com/go-redis/redis"
	"log"
)

//单点redis
var RedisCli *redis.Client

func Init(opt *IOptions) {
	//初始默认配置
	opt.Init()
	RedisCli = redis.NewClient(&redis.Options{
		Addr:        fmt.Sprintf("%s:%d", opt.Host, opt.Port),
		Password:    opt.Password,
		DB:          opt.DB,
		MaxRetries:  opt.MaxRetry,
		DialTimeout: opt.DialTimeOut,
		MaxConnAge:  opt.MaxConnAge,
	})
	if RedisCli.Ping().Err() != nil {
		log.Fatal("Redis初始化失败.")
	}
}
