package main

import (
	"HelloMyWorld/common/iredis"
	"HelloMyWorld/config"
	"testing"
)

func TestRedisInit(t *testing.T) {
	config.Init()
	iredis.Init(&iredis.IOptions{
		Host:        "127.0.0.1",
		Port:        6379,
		DB:          0,
		Password:    "123456",
		MaxRetry:    5,
		DialTimeOut: 10000,//秒
		MaxConnAge:  3,//分钟
	})
	iredis.RedisCli.Set("Hai","GeGe",0)
}
