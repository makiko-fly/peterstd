package peterstd

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"github.com/sirupsen/logrus"
	"gitlab.wallstcn.com/spider/peterstd/logging"
	"gitlab.wallstcn.com/spider/peterstd/util"
)

const (
	ContextRequestID = "request-id"
	ContextMetadata  = "context-metadata"
	ContextLogger    = "logger"
)

type (
	LoggingConfig = logging.LoggingConfig
	ILogger       = *logrus.Entry
)

func Logger() ILogger {
	return logging.GetLogger()
}

// 将 context 的 metadata 中的字段放入logger中
func WithMetadata(c context.Context) ILogger {
	logger := Logger()
	md := c.Value(ContextMetadata).(map[string]string)
	for k, v := range md {
		logger = logger.WithField(k, v)
	}
	return logger
}

func With(key string, value interface{}) ILogger {
	return Logger().WithField(key, value)
}

func WithField(key string, value interface{}) ILogger {
	return Logger().WithField(key, value)
}

func WithDebugFlag(value interface{}) ILogger {
	return Logger().WithField("debug_flag", value)
}

func WithTraceID(value string) ILogger {
	return Logger().WithField("dev-trace-id", value)
}

func WithLineNum() ILogger {
	lineInfo := util.GetLineNum()
	if lineInfo == "" {
		return Logger()
	}
	return Logger().WithField("caller", lineInfo)
}

func WithElapsedTime(start time.Time) ILogger {
	return Logger().WithField("elapsed_time", time.Since(start).Seconds()*1000)
}

func WithError(err error) ILogger {
	return Logger().WithError(err)
}

func Print(v ...interface{}) {
	Logger().Print(v...)
}

func Debug(v ...interface{}) {
	Logger().Debug(v...)
}

func Info(v ...interface{}) {
	Logger().Info(v...)
}

func Warn(v ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	caller := fmt.Sprintf("%s: %d", file, line)
	Logger().WithField("caller", caller).Warn(v...)
}

func Error(v ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	caller := fmt.Sprintf("%s: %d", file, line)
	Logger().WithField("caller", caller).Error(v...)
}

func Fatal(v ...interface{}) {
	Logger().Fatal(v...)
}

func Panic(v ...interface{}) {
	Logger().Panic(v...)
}

func Println(v ...interface{}) {
	Logger().Println(v...)
}

func Debugln(v ...interface{}) {
	Logger().Debugln(v...)
}

func Infoln(v ...interface{}) {
	Logger().Infoln(v...)
}

func Warnln(v ...interface{}) {
	Logger().Warnln(v...)
}

func Errorln(v ...interface{}) {
	Logger().Errorln(v...)
}

func Fatalln(v ...interface{}) {
	Logger().Fatalln(v...)
}

func Panicln(v ...interface{}) {
	Logger().Panicln(v...)
}

func Printf(fmt string, v ...interface{}) {
	Logger().Printf(fmt, v...)
}

func Debugf(fmt string, v ...interface{}) {
	Logger().Debugf(fmt, v...)
}

func Infof(fmt string, v ...interface{}) {
	Logger().Infof(fmt, v...)
}

func Warnf(fmt string, v ...interface{}) {
	Logger().Warnf(fmt, v...)
}

func Errorf(fmt string, v ...interface{}) {
	Logger().Errorf(fmt, v...)
}

func Fatalf(fmt string, v ...interface{}) {
	Logger().Fatalf(fmt, v...)
}

func Panicf(fmt string, v ...interface{}) {
	Logger().Panicf(fmt, v...)
}
