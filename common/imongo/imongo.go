package imongo

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

var DB *mongo.Database

type IOptions struct {
	Host       string
	Port       int
	DB         string
	User       string
	Password   string
	AuthSource string
	TimeOut    int
}

func Init(opt *IOptions) {
	//配置参数初始化
	opt.Init()
	//链接参数
	url := fmt.Sprintf("mongodb://%s:%s@%s:%d/%s?authSource=%s", opt.User, opt.Password, opt.Host, opt.Port, opt.DB, opt.AuthSource)

	mongoOpt := new(options.ClientOptions)
	mongoOpt.ApplyURI(url)
	//设置超时
	cxt, _ := context.WithTimeout(context.Background(), time.Second*time.Duration(opt.TimeOut))
	client, err := mongo.Connect(cxt, mongoOpt)

	if err != nil {
		log.Fatalf("Connect Mongo Err:%v.\n", err)
	}

	DB = client.Database(opt.DB)
}

func (p *IOptions) Init() {
	if len(p.Host) == 0 {
		p.Host = "127.0.0.1"
	}
	if p.Port == 0 {
		p.Port = 27017
	}
	if len(p.User) == 0 {
		p.User = "root"
	}
	if len(p.Password) == 0 {
		p.Password = "123456"
	}
	if p.TimeOut == 0 {
		//默认三秒超时
		p.TimeOut = 3
	}
}
