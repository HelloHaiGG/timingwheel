package main

import (
	"HelloMyWorld/common/iredis"
	"HelloMyWorld/config"
	"testing"
)

func TestRedisInit(t *testing.T) {
	config.Init()
	iredis.Init(&iredis.IOptions{
		Host:        "",
		Port:        0,
		DB:          0,
		Password:    "",
		MaxRetry:    0,
		DialTimeOut: 0,
		MaxConnAge:  0,
	})
}
