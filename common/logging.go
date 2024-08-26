package common

import (
	"os"
	"time"
	"github.com/sirupsen/logrus"
)

var (
	errorLogger *logrus.Logger
	todoLogger  *logrus.Logger
)

func init() {
	// error.log ロガーの初期化
	errorLogger = logrus.New()
	errorLogger.SetFormatter(&logrus.JSONFormatter{})
	file, err := os.OpenFile("logs/error.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		logrus.Fatalf("Failed to open error log file: %v", err)
	}
	errorLogger.SetOutput(file)

	// todo.log ロガーの初期化
	todoLogger = logrus.New()
	todoLogger.SetFormatter(&logrus.JSONFormatter{})
	file, err = os.OpenFile("logs/todo.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		logrus.Fatalf("Failed to open todo log file: %v", err)
	}
	todoLogger.SetOutput(file)
}

func LogError(err error) {
	errorLogger.WithFields(logrus.Fields{
		"timestamp": time.Now().Format(time.RFC3339),
		"level":     "ERROR",
		"message":   err.Error(),
	}).Error()
}

type LogLevel string

const (
	INFO  LogLevel = "INFO"
	WARN  LogLevel = "WARN"
	ERROR LogLevel = "ERROR"
)

func LogTodo(level LogLevel, message string) {
	todoLogger.WithFields(logrus.Fields{
		"timestamp": time.Now().Format(time.RFC3339),
		"level":     level,
		"message":   message,
	}).Info()
}
