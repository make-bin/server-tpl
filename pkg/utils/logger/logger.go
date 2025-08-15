package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

// InitLogger 初始化日志
func InitLogger(level, format, output string) {
	log = logrus.New()

	// 设置日志级别
	logLevel, err := logrus.ParseLevel(level)
	if err != nil {
		logLevel = logrus.InfoLevel
	}
	log.SetLevel(logLevel)

	// 设置日志格式
	if format == "json" {
		log.SetFormatter(&logrus.JSONFormatter{})
	} else {
		log.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
	}

	// 设置输出
	if output == "stdout" {
		log.SetOutput(os.Stdout)
	} else if output == "stderr" {
		log.SetOutput(os.Stderr)
	} else {
		file, err := os.OpenFile(output, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.SetOutput(os.Stdout)
		} else {
			log.SetOutput(file)
		}
	}
}

// GetLogger 获取日志实例
func GetLogger() *logrus.Logger {
	if log == nil {
		InitLogger("info", "json", "stdout")
	}
	return log
}

// WithField 添加字段到日志
func WithField(key string, value interface{}) *logrus.Entry {
	return GetLogger().WithField(key, value)
}

// WithFields 添加多个字段到日志
func WithFields(fields logrus.Fields) *logrus.Entry {
	return GetLogger().WithFields(fields)
}

// Debug 调试日志
func Debug(args ...interface{}) {
	GetLogger().Debug(args...)
}

// Info 信息日志
func Info(args ...interface{}) {
	GetLogger().Info(args...)
}

// Warn 警告日志
func Warn(args ...interface{}) {
	GetLogger().Warn(args...)
}

// Error 错误日志
func Error(args ...interface{}) {
	GetLogger().Error(args...)
}

// Fatal 致命错误日志
func Fatal(args ...interface{}) {
	GetLogger().Fatal(args...)
}

// Debugf 格式化调试日志
func Debugf(format string, args ...interface{}) {
	GetLogger().Debugf(format, args...)
}

// Infof 格式化信息日志
func Infof(format string, args ...interface{}) {
	GetLogger().Infof(format, args...)
}

// Warnf 格式化警告日志
func Warnf(format string, args ...interface{}) {
	GetLogger().Warnf(format, args...)
}

// Errorf 格式化错误日志
func Errorf(format string, args ...interface{}) {
	GetLogger().Errorf(format, args...)
}

// Fatalf 格式化致命错误日志
func Fatalf(format string, args ...interface{}) {
	GetLogger().Fatalf(format, args...)
}
