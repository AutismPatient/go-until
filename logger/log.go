package logger

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"time"
)

var (
	DefaultLogouts = &logrus.Logger{}
)

func init() {
	DefaultLogouts = NewLog(Parameter{
		Level:        5,
		ReportCaller: true,
		Fields:       nil,
		Hook:         nil,
		IO:           os.Stdout,
		RegisterExitHandler: func() {

		},
		DeferExitHandler: func() {

		},
	})
}

/**
参数
*/
type Parameter struct {
	Level               logrus.Level  // 级别
	ReportCaller        bool          // 如果您希望将调用方法添加为字段
	Fields              logrus.Fields // 自定义字段map
	Hook                logrus.Hook   // 您可以添加用于日志记录级别的挂钩。例如，将错误发送到上的异常跟踪服务Error，Fatal并将Panic信息发送到StatsD或同时记录到多个位置，例如syslog。
	IO                  io.Writer     // io
	RegisterExitHandler func()        // 程序异常后即将退出事件
	DeferExitHandler    func()        //
}

/**
返回 logrus 实例
*/
func NewLog(parameter Parameter) *logrus.Logger {
	DefaultLogouts = logrus.New()

	if parameter.Level == 0 {
		parameter.Level = 5
	}

	DefaultLogouts.SetLevel(parameter.Level)

	DefaultLogouts.SetReportCaller(parameter.ReportCaller)
	DefaultLogouts.WithFields(parameter.Fields)
	if parameter.Hook != nil {
		DefaultLogouts.AddHook(parameter.Hook)
	}
	if parameter.IO != nil {
		DefaultLogouts.SetOutput(parameter.IO)
	}

	logrus.RegisterExitHandler(parameter.RegisterExitHandler)
	logrus.DeferExitHandler(parameter.DeferExitHandler)

	return DefaultLogouts
}

/**
默认格式化输出 todo 2021年2月6日22:42:06
*/
type LogoutsFormatter struct {
}

func (receiver *LogoutsFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}
	b.WriteString(fmt.Sprintf("%s %s", time.Now().Format("2006-01-02 15:04:05"), entry.Message))
	return b.Bytes(), nil
}

/**
日志记录器
来自： https://github.com/sirupsen/logrus
*/
