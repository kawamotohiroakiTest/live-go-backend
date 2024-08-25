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
		fmt.Println("Error loading .env file")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r := mux.NewRouter()
	r.HandleFunc("/api/v1/health", common.HealthHandler)
	r.HandleFunc("/api/v1/todo/{id}", common.TodoHandler)

	fmt.Println("Starting server on! :" + port)
	if err := http.ListenAndServe(":"+port, common.EnableCors(r)); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
