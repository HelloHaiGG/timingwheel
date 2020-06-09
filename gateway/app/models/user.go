package models

import "github.com/jinzhu/gorm"

//用户发送短信
type SendSMSParams struct {
	Phone string `json:"phone"`
	ReqId string `json:"req_id"`
	Code  string `json:"code"`
}

//用户短信验证码验证
type CheckSMSCodeParams struct {

}

//用户结构
type User struct {
	gorm.Model
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

//用户注册
type RegisterParams struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
	VerifyId string `json:"verify_id"`
	Code     string `json:"code"`
}
