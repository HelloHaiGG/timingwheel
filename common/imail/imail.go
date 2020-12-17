package imail

import (
	"github.com/jordan-wright/email"
	"log"
	"net/smtp"
)

type IOption struct {
	Service  string `json:"service"`   //"邮件服务商"
	AuthCode string `json:"auth_code"` //授权码
	Sender   string `json:"sender"`    //发送者
}

var Ins *email.Pool

func Init(ops *IOption) {
	if ops.Service == "" || ops.AuthCode == "" {
		log.Fatal("Init Email Server Fail. ")
		return
	}
	pool, err := email.NewPool(ops.Service+":25", 10, smtp.PlainAuth("", ops.Sender, ops.AuthCode, ops.Service))
	if err != nil {
		log.Fatal("New Email Server Pool Fail. err: ", err)
		return
	}
	Ins = pool
}
