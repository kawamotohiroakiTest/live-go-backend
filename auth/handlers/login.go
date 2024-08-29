package handlers

import (
	"encoding/json"
	"live/auth/models"
	"live/common"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type LoginCredentials struct {
	Mail string `json:"mail" validate:"required,email"`
	Pass string `json:"pass" validate:"required"`
}

var jwtKey = []byte("your_secret_key")

type Claims struct {
	UserID uint   `json:"user_id"`
	Mail   string `json:"mail"`
	jwt.StandardClaims
}

func Login(w http.ResponseWriter, r *http.Request) {
	var creds LoginCredentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		common.LogUser(common.ERROR, err.Error())
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// バリデーションの実行
	validate := validator.New()
	err = validate.Struct(creds)
	if err != nil {
		common.LogUser(common.ERROR, err.Error())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Validation failed", "details": err.Error()})
		return
	}

	var user models.User
	if err := common.DB.Where("mail = ?", creds.Mail).First(&user).Error; err != nil {
		common.LogUser(common.ERROR, "User not found: "+creds.Mail)
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// パスワードの比較
	if err := bcrypt.CompareHashAndPassword([]byte(user.Pass), []byte(creds.Pass)); err != nil {
		common.LogUser(common.ERROR, "Invalid password for user: "+creds.Mail)
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// JWTトークンの作成
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID: user.ID,
		Mail:   user.Mail,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		common.LogUser(common.ERROR, "Failed to generate JWT: "+err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// トークンをレスポンスとして返す
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}
