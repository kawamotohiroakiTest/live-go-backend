package common

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB() {
	// MySQLコンテナへの接続情報
	dsn := "liveuser:livepass@tcp(db:3306)/live"
	var err error

	// DBオブジェクトの初期化
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to MySQL: %v", err)
	}

	// DB接続のテスト
	err = DB.Ping()
	if err != nil {
		log.Fatalf("Failed to ping MySQL: %v", err)
	}

	fmt.Println("Successfully connected to MySQL")
}
