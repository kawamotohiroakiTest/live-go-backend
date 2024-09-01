package common

import (
	"fmt"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() {
	user := os.Getenv("MYSQL_USER")
	password := os.Getenv("MYSQL_PASSWORD")
	host := os.Getenv("MYSQL_HOST")
	port := os.Getenv("MYSQL_PORT")
	database := os.Getenv("MYSQL_DATABASE")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, host, port, database)

	var err error
	for i := 0; i < 100; i++ {
		DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent), // ログの詳細度を設定
		})
		if err == nil {
			sqlDB, err := DB.DB()
			if err != nil {
				fmt.Println("Failed to get sqlDB from GORM DB: ", err)
				continue
			}
			err = sqlDB.Ping()
			if err == nil {
				fmt.Println("Successfully connected to MySQL with GORM")
				return
			}
		}
		fmt.Println("Failed to connect to MySQL with GORM. Retrying...")
		time.Sleep(10 * time.Second)
	}
	if err != nil {
		fmt.Printf("Error: Could not connect to MySQL after multiple attempts: %v\n", err)
	}
}
