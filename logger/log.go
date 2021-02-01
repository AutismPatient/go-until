package logger

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
)

var (
	DefaultLogouts = &logrus.Logger{}
)

func init() {

}

/**
参数
*/
type Parameter struct {
	Formatter           LogoutsFormatter // 自定义格式化
	Level               logrus.Level     // 级别
	ReportCaller        bool             // 如果您希望将调用方法添加为字段
	Fields              logrus.Fields    // 自定义字段map
	Hook                logrus.Hook      // 您可以添加用于日志记录级别的挂钩。例如，将错误发送到上的异常跟踪服务Error，Fatal并将Panic信息发送到StatsD或同时记录到多个位置，例如syslog。
	IO                  io.Writer        // io
	RegisterExitHandler func()           // 程序异常后即将退出事件
	DeferExitHandler    func()           //
}

/**
返回 logrus 实例
*/
func newLog(parameter Parameter) *logrus.Logger {
	DefaultLogouts = logrus.New()
	DefaultLogouts.SetFormatter(&parameter.Formatter)
	DefaultLogouts.SetLevel(parameter.Level)
	DefaultLogouts.SetReportCaller(parameter.ReportCaller)
	DefaultLogouts.WithFields(parameter.Fields)
	DefaultLogouts.AddHook(parameter.Hook)
	DefaultLogouts.SetOutput(parameter.IO)

	logrus.RegisterExitHandler(parameter.RegisterExitHandler)
	logrus.DeferExitHandler(parameter.DeferExitHandler)

	return DefaultLogouts
}

/**
格式化输出
*/
type LogoutsFormatter struct {
}

func (receiver *LogoutsFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	serialized, err := json.Marshal(entry.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal fields to JSON, %w", err)
	}
	return append(serialized, '\n'), nil
}

/**
日志记录器
来自： https://github.com/sirupsen/logrus
*/
