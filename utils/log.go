package main

import (
	"github.com/sirupsen/logrus"
	"os"
)

var Log *logrus.Logger

// InitLog 日志库
func InitLog() {
	Log = logrus.New()
	Log.Out = os.Stdout
	Log.Formatter = &logrus.JSONFormatter{}
}
