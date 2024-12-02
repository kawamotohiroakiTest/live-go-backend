package main

import (
	"flag"
	"fmt"
	"live/auth"
	"live/common"
	"live/db"
	"live/videohub"
	"live/videoupload"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	// フラグのパース
	flag.Parse()

	err := godotenv.Load()
	if err != nil {
		common.LogError(fmt.Errorf("Error loading .env file: %v", err))
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	common.InitDB()

	// コマンドラインフラグのチェック
	if len(flag.Args()) > 0 && flag.Args()[0] == "migrate" {
		db.RunMigration()
		return
	}

	// マイグレーションの実行
	common.LogTodo(common.INFO, "Running database migrations...")
	db.RunMigration()

	r := mux.NewRouter()

	auth.RegisterRoutes(r)
	videoupload.RegisterRoutes(r)
	videohub.RegisterRoutes(r)

	r.HandleFunc("/api/v1/health", common.HealthHandler)
	r.HandleFunc("/api/v1/todo/{id}", common.TodoHandler)

	common.LogTodo(common.INFO, "Starting server on port!: "+port)
	if err := http.ListenAndServe(":"+port, common.EnableCors(r)); err != nil {
		common.LogError(fmt.Errorf("Error starting server: %v", err))
	}
}
