package common

import (
	"os"
	"time"
	"github.com/sirupsen/logrus"
	"io"
)

var (
	errorLogger *logrus.Logger
	todoLogger  *logrus.Logger
)

func init() {
	// logsディレクトリの存在を確認し、なければ作成
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		err := os.Mkdir("logs", 0755) 
		if err != nil {
			logrus.Fatalf("Failed to create logs directory: %v", err)
		}
	}

	// error.log ロガーの初期化
	errorLogger = logrus.New()
	errorLogger.SetFormatter(&logrus.JSONFormatter{})
	errorLogFile, err := os.OpenFile("logs/error.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		logrus.Fatalf("Failed to open error log file: %v", err)
	}

	// stdout と error.log に同時にログを出力するためにマルチライターを設定
	errorMultiWriter := io.MultiWriter(os.Stderr, errorLogFile)
	errorLogger.SetOutput(errorMultiWriter)

	// todo.log ロガーの初期化
	todoLogger = logrus.New()
	todoLogger.SetFormatter(&logrus.JSONFormatter{})
	todoLogFile, err := os.OpenFile("logs/todo.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		logrus.Fatalf("Failed to open todo log file: %v", err)
	}

	// stdout と todo.log に同時にログを出力するためにマルチライターを設定
	todoMultiWriter := io.MultiWriter(os.Stdout, todoLogFile)
	todoLogger.SetOutput(todoMultiWriter)
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
