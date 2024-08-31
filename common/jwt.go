package common

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

// JWT秘密鍵の定義
var JwtKey = []byte("your_secret_key")

type Claims struct {
	UserID uint   `json:"user_id"`
	Mail   string `json:"mail"`
	jwt.StandardClaims
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return JwtKey, nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		if !token.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// 次のハンドラにクレーム情報を渡す
		ctx := context.WithValue(r.Context(), "claims", claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserIDFromContext(ctx context.Context) (uint, error) {
	// claims 変数の内容をログに出力
	claims, ok := ctx.Value("claims").(*Claims)
	// LogVideoUploadError(fmt.Errorf("claims: %v, ok: %v", claims, ok))

	if !ok || claims == nil {
		err := fmt.Errorf("Failed to get claims from context")
		LogVideoUploadError(err)
		return 0, err
	}

	// 最終的に返す UserID をログに出力
	// LogVideoUploadError(fmt.Errorf("UserID: %d", claims.UserID))

	return claims.UserID, nil
}
