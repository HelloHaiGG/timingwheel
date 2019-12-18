package ilogger

import (
	"fmt"
	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
	"os"
)

//是否将日志写入到文件
var ToFile bool

type ILogger struct {
	instance *logrus.Logger
}

var sev string

func Init(server string) *ILogger {
	if server != "" {
		sev = server
	} else {
		panic("ilogger server is nil!")
	}
	iLogger := &ILogger{instance: logrus.New()}
	if ToFile {
		//设置为json格式
		iLogger.instance.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
			FieldMap: logrus.FieldMap{
				"time": "when",
				//"file":  "where",
				"msg":   "what",
				"level": "level",
			}})
		path, _ := os.Getwd()
		iLogger.instance.SetOutput(&lumberjack.Logger{
			Filename: path + fmt.Sprintf("/logs/%s.txt", server),
			MaxSize:  200, //最大 分割
		})

	} else {
		iLogger.instance.SetFormatter(&logrus.TextFormatter{
			ForceColors:     true,
			DisableSorting:  false,
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
			FieldMap: logrus.FieldMap{
				"time": "when",
				//"file":  "where",
				"msg":   "what",
				"level": "level",
			}})
		iLogger.instance.SetOutput(os.Stdout)
	}
	//调用者
	//iLogger.instance.SetReportCaller(true)
	iLogger.instance.AddHook(&IHooks{
		levels: logrus.AllLevels,
		Field:  "where",
		Skip:   9,
	})
	return iLogger
}

func (p *ILogger) Info(args ...interface{}) {
	p.instance.SetLevel(logrus.InfoLevel)
	p.instance.WithField("who", sev).Info(args...)
}

func (p *ILogger) Error(args ...interface{}) {
	p.instance.SetLevel(logrus.ErrorLevel)
	p.instance.WithField("who", sev).Error(args...)
}

func (p *ILogger) Debug(args ...interface{}) {
	p.instance.SetLevel(logrus.DebugLevel)
	p.instance.WithField("who", sev).Debug(args...)
}

func (p *ILogger) Warn(args ...interface{}) {
	p.instance.SetLevel(logrus.WarnLevel)
	p.instance.WithField("who", sev).Warn(args...)
}

//fatal--> 打印 log 之后会调用 os.Exit(1)
func (p *ILogger) Exit(args ...interface{}) {
	p.instance.SetLevel(logrus.FatalLevel)
	p.instance.WithField("who", sev).Fatal(args...)
}

//打印log 之后会 panic
func (p *ILogger) Panic(args ...interface{}) {
	p.instance.SetLevel(logrus.PanicLevel)
	p.instance.WithField("who", sev).Fatal(args...)
}
