package main

import (
	"HelloMyWorld/common/igorm"
	"HelloMyWorld/config"
	"fmt"
	"testing"
)

type Db struct {
	Database string
}

func TestGormInit(t *testing.T) {
	config.Init()
	igorm.Init(&igorm.IOption{
		Host:     config.APPConfig.Mysql.Host,
		Port:     0,
		User:     "root",
		Password: "123456",
		DB:       "seekspider",
		IsDebug:  true,
	})
	fmt.Println(config.APPConfig.Mysql)
	if err := igorm.DbClient.DB().Ping(); err != nil {
		t.Log(err)
	}
}
