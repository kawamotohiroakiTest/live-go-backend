package handlers

import (
	"encoding/json"
	"live/auth/models"
	"live/common"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type Credentials struct {
	Name     string `json:"name"`
	Mail     string `json:"mail"`
	Password string `json:"password"`
}

func Register(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		common.LogUser(common.ERROR, err.Error())
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
	if err != nil {
		common.LogUser(common.ERROR, err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	tx := common.DB.Begin()
	if tx.Error != nil {
		common.LogUser(common.ERROR, tx.Error.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	user := models.User{Name: creds.Name, Mail: creds.Mail, Pass: string(hashedPassword)}

	if err := user.Validate(); err != nil {
		common.LogUser(common.ERROR, err.Error())
		tx.Rollback()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Validation failed", "details": err.Error()})
		return
	}

	if err := tx.Create(&user).Error; err != nil {
		common.LogUser(common.ERROR, err.Error())
		tx.Rollback()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to register user"})
		return
	}

	if err := tx.Commit().Error; err != nil {
		common.LogUser(common.ERROR, err.Error())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to register user"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
}
