package common

import (
	"net/http"
	"strconv"
	"fmt"
	"encoding/json"
	"github.com/gorilla/mux"
)

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{
		"status": "OK",
	}
	json.NewEncoder(w).Encode(response)
}

func TodoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		LogError(err)  // エラーログを書き込む
		return
	}

	var foundTodo *Todo
	for _, todo := range Todos {
		if todo.ID == id {
			foundTodo = &todo
			break
		}
	}

	if foundTodo == nil {
		errMsg := fmt.Errorf("Todo not found with ID %d", id)
		http.Error(w, errMsg.Error(), http.StatusNotFound)
		LogTodo(WARN, fmt.Sprintf("Todo not found for user input ID %d", id))  // ユーザーの入力に関連するエラーをログに記録
		return
	}

	LogTodo(INFO, fmt.Sprintf("Successfully retrieved Todo with ID %d", id))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(foundTodo)
}
