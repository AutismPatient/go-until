package logger

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
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
	Formatter    LogoutsFormatter // 自定义格式化
	Level        logrus.Level     // 级别
	ReportCaller bool             // 如果您希望将调用方法添加为字段
	Fields       logrus.Fields    // 自定义字段map
}

/**
go日志封装
来自：log
*/
func newLog(parameter Parameter) *logrus.Logger {
	DefaultLogouts = logrus.New()
	DefaultLogouts.SetFormatter(&parameter.Formatter)
	DefaultLogouts.SetLevel(parameter.Level)
	DefaultLogouts.SetReportCaller(parameter.ReportCaller)
	DefaultLogouts.WithFields(parameter.Fields)
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
