package common

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB() {

	user := os.Getenv("MYSQL_USER")
	password := os.Getenv("MYSQL_PASSWORD")
	host := os.Getenv("MYSQL_HOST")
	port := os.Getenv("MYSQL_PORT")
	database := os.Getenv("MYSQL_DATABASE")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, port, database)

	var err error
	for i := 0; i < 5; i++ {
		DB, err = sql.Open("mysql", dsn)
		if err == nil && DB.Ping() == nil {
			fmt.Println("Successfully connected to MySQL")
			return
		}
		fmt.Println("Failed to connect to MySQL. Retrying...")
		time.Sleep(5 * time.Second)
	}
}
