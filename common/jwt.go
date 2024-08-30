package common

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

// JWT秘密鍵の定義
var JwtKey = []byte("your_secret_key")

// JWTクレームの構造体
type Claims struct {
	UserID uint   `json:"user_id"`
	Mail   string `json:"mail"`
	jwt.StandardClaims
}

// jWT認証をチェックする
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

// jwtからuser_idを取得する
func GetUserIDFromContext(ctx context.Context) (uint, error) {
	claims, ok := ctx.Value("claims").(*Claims)
	if !ok || claims == nil {
		return 0, jwt.ErrSignatureInvalid
	}
	return claims.UserID, nil
}
