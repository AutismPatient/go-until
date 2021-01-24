package logger

import (
	"github.com/sirupsen/logrus"
	"log"
)

var (
	Logger  = &log.Logger{}
	Logouts = &logrus.Logger{}
)

func init() {

}

/**
go日志封装
来自：log
*/
func newLog() {
	Logouts = logrus.New()
}

/**
日志记录器
来自： https://github.com/sirupsen/logrus
*/
