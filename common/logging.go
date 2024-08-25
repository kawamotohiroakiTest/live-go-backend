package common

import (
	"log"
	"os"
	"time"
)

func LogError(err error) {
	// error.logファイルを開く、または作成
	file, fileErr := os.OpenFile("logs/error.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if fileErr != nil {
		log.Printf("Failed to open error log file: %v\n", fileErr)
		return
	}
	defer file.Close()

	// ロガーを作成
	logger := log.New(file, "", log.LstdFlags)
	logger.Printf("[%s] ERROR: %v\n", time.Now().Format(time.RFC3339), err)
}

type LogLevel string

const (
	INFO  LogLevel = "INFO"
	WARN  LogLevel = "WARN"
	ERROR LogLevel = "ERROR"
)

func LogTodo(level LogLevel, message string) {
	// todo.logファイルを開く、または作成
	file, fileErr := os.OpenFile("logs/todo.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if fileErr != nil {
		log.Printf("Failed to open todo log file: %v\n", fileErr)
		return
	}
	defer file.Close()

	// ロガーを作成
	logger := log.New(file, "", log.LstdFlags)
	logger.Printf("[%s] %s: %s\n", time.Now().Format(time.RFC3339), level, message)
}