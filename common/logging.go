package common

import (
	"io"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	errorLogger       *logrus.Logger
	todoLogger        *logrus.Logger
	userLogger        *logrus.Logger
	videouploadLogger *logrus.Logger
)

func init() {
	// logsディレクトリの存在を確認し、なければ作成
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		err := os.Mkdir("logs", 0755)
		if err != nil {
			logrus.Fatalf("Failed to create logs directory: %v", err)
		}
	}

	// 共通の標準出力設定
	stdOut := os.Stdout

	// error.log ロガーの初期化
	errorLogger = logrus.New()
	errorLogger.SetFormatter(&logrus.JSONFormatter{
		DisableHTMLEscape: true,
	})
	errorLogFile, err := os.OpenFile("logs/error.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		logrus.Fatalf("Failed to open error log file: %v", err)
	}
	// 標準出力とファイルの両方に書き込み
	errorLogger.SetOutput(io.MultiWriter(stdOut, errorLogFile))

	// todo.log ロガーの初期化
	todoLogger = logrus.New()
	todoLogger.SetFormatter(&logrus.JSONFormatter{
		DisableHTMLEscape: true,
	})
	todoLogFile, err := os.OpenFile("logs/todo.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		logrus.Fatalf("Failed to open todo log file: %v", err)
	}
	// 標準出力とファイルの両方に書き込み
	todoLogger.SetOutput(io.MultiWriter(stdOut, todoLogFile))

	// user.log ロガーの初期化
	userLogger = logrus.New()
	userLogger.SetFormatter(&logrus.JSONFormatter{
		DisableHTMLEscape: true,
	})
	userLogFile, err := os.OpenFile("logs/user.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		logrus.Fatalf("Failed to open user log file: %v", err)
	}
	// 標準出力とファイルの両方に書き込み
	userLogger.SetOutput(io.MultiWriter(stdOut, userLogFile))

	// videoupload.log ロガーの初期化
	videouploadLogger = logrus.New()
	videouploadLogger.SetFormatter(&logrus.JSONFormatter{
		DisableHTMLEscape: true,
	})
	videouploadLogFile, err := os.OpenFile("logs/videoupload.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		logrus.Fatalf("Failed to open videoupload log file: %v", err)
	}
	// 標準出力とファイルの両方に書き込み
	videouploadLogger.SetOutput(io.MultiWriter(stdOut, videouploadLogFile))
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

func LogUser(level LogLevel, message string) {
	userLogger.WithFields(logrus.Fields{
		"timestamp": time.Now().Format(time.RFC3339),
		"level":     level,
		"message":   message,
	}).Info()
}

func LogVideoUploadError(err error) {
	videouploadLogger.WithFields(logrus.Fields{
		"timestamp": time.Now().Format(time.RFC3339),
		"level":     "ERROR",
		"message":   err.Error(),
	}).Error()
}
