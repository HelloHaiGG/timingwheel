package iRedis

import "time"

//redis可动态配置的选项
type IOptions struct {
	Host        string
	Port        int
	DB          int
	Password    string
	MaxRetry    int
	DialTimeOut time.Duration
	//空闲连接的保活时常
	MaxConnAge time.Duration
}

//初始化默认值
func (p *IOptions) Init() {
	if p.Host == "" {
		p.Host = "127.0.0.1"
	}
	if p.Port == 0 {
		p.Port = 6379
	}
	if p.DB == 0 {
		p.DB = 0
	}
	if p.Password == "" {
		p.Password = "root"
	}
	if p.MaxRetry == 0 {
		p.MaxRetry = 5
	}
	if p.DialTimeOut == 0 {
		p.DialTimeOut = 10
	}
	//五分钟保活
	if p.MaxConnAge == 0 {
		p.MaxConnAge = time.Minute * 5
	}
}
