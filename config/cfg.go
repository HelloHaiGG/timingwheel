package config

import (
	"HelloMyWorld/utils"
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
)

var APPConfig AppCfg

//解析配置文件
func Init() {
	pwd, _ := os.Getwd()
	path := fmt.Sprintf("%s/%s", pwd, "config.yaml")
	if !utils.IsExist(path) {
		panic(fmt.Sprintf("%s does not exist.",path))
	}
	if b, err := utils.HandFile(path); err != nil {
		panic(fmt.Sprintf("%s loading error:%v",path,err))
	} else {
		if err = yaml.Unmarshal(b, &APPConfig); err != nil {
			return
			panic(fmt.Sprintf("Unmarshal appconfig error:%v",err))
		}
	}
}
