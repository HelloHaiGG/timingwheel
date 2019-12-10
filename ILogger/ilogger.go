package ILogger

import (
	"fmt"
	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
	"os"
)

type ILogger struct {
	instance *logrus.Logger
}

var sev string

func Init(server string) *ILogger {
	if server != "" {
		sev = server
	} else {
		panic("ILogger server is nil!")
	}
	iLogger := &ILogger{instance: logrus.New()}
	//设置输出位置
	// 本地调试 running = ""，打印到标准输出
	// 线上调试 running = 1，输出到日志文件
	if os.Getenv("RUNNING") == "1" {
		//设置为json格式
		iLogger.instance.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
			FieldMap: logrus.FieldMap{
				"time":  "when",
				"file":  "where",
				"msg":   "what",
				"level": "level",
			}})
		path, _ := os.Getwd()
		iLogger.instance.SetOutput(&lumberjack.Logger{
			Filename: path + fmt.Sprintf("\\logs\\%s.txt", server),
			MaxSize:  200, //最大 分割
		})
		fmt.Println(path)
	} else {
		iLogger.instance.SetFormatter(&logrus.TextFormatter{
			ForceColors:     true,
			DisableSorting:  true,
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
			FieldMap: logrus.FieldMap{
				"time":  "when",
				"file":  "where",
				"msg":   "what",
				"level": "level",
			}})
		iLogger.instance.SetOutput(os.Stdout)
	}
	iLogger.instance.SetReportCaller(true)
	return iLogger
}

func (p *ILogger) Info(args ...interface{}) {
	p.instance.WithField("who", sev).Info(args...)
}

func (p *ILogger) Error(args ...interface{}) {
	p.instance.WithField("who", sev).Error(args...)
}

func (p *ILogger) Debug(args ...interface{}) {
	p.instance.WithField("who", sev).Debug(args...)
}

func (p *ILogger) Warn(args ...interface{}) {
	p.instance.WithField("who", sev).Warn(args...)
}

//fatal--> 打印 log 之后会调用 os.Exit(1)
func (p *ILogger) Exit(args ...interface{}) {
	p.instance.WithField("who", sev).Fatal(args...)
}

//打印log 之后会 panic
func (p *ILogger) Panic(args ...interface{}) {
	p.instance.WithField("who", sev).Fatal(args...)
}
