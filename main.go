package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"encoding/json"

	"github.com/joho/godotenv"
	"github.com/gorilla/mux"
)

// Todo構造体を定義
type Todo struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

var todos = []Todo{
	{ID: 1, Title: "Sample Todo 1", Completed: false},
	{ID: 2, Title: "Sample Todo 2", Completed: true},
	{ID: 3, Title: "Sample Todo 3", Completed: false},
}

// CORS対応のためのミドルウェア
func enableCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// CORSヘッダーを追加
		allowedOrigin := os.Getenv("API_ALLOWED_ORIGIN")
		w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// プリフライトリクエスト（OPTIONSメソッド）への対応
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// 次のハンドラを呼び出す
		next.ServeHTTP(w, r)
	})
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    response := map[string]string{
        "status": "OK",
    }
    json.NewEncoder(w).Encode(response)
}

// ハンドラ関数: /api/v1/todo/{id} に対応するTodoを返す
func todoHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }

    // ダミーデータを作成
    response := map[string]interface{}{
        "id":        id,
        "title":     fmt.Sprintf("Sample Todo %d", id),
        "completed": id%2 == 0, // 偶数のIDは完了済み、奇数のIDは未完了
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}


func main() {
    err := godotenv.Load()
    if err != nil {
        fmt.Println("Error loading .env file")
    }

    // 環境変数からポート番号を取得
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080" // デフォルトポート
    }

    // ルーターを作成
    r := mux.NewRouter()
	r.HandleFunc("/api/v1/health", healthHandler)
    r.HandleFunc("/api/v1/todo/{id}", todoHandler)

    // CORS対応ミドルウェアを適用してサーバーを起動
    fmt.Println("Starting server on! :" + port)
    if err := http.ListenAndServe(":"+port, enableCors(r)); err != nil {
        fmt.Println("Error starting server:", err)
    }
}
