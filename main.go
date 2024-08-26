package main

import (
	"fmt"
	"net/http"
	"os"
	"github.com/joho/godotenv"
	"github.com/gorilla/mux"
	"live/common"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		common.LogError(fmt.Errorf("Error loading .env file: %v", err))
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r := mux.NewRouter()
	r.HandleFunc("/api/v1/health", common.HealthHandler)
	r.HandleFunc("/api/v1/todo/{id}", common.TodoHandler)

	common.LogTodo(common.INFO, "Starting server on port: "+port)
	if err := http.ListenAndServe(":"+port, common.EnableCors(r)); err != nil {
		common.LogError(fmt.Errorf("Error starting server: %v", err))
	}
}
