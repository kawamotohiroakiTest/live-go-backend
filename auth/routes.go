package auth

import (
	"live/auth/handlers"
	"live/common"
	"net/http"

	"github.com/gorilla/mux"
)

func RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/api/v1/users/register", handlers.Register).Methods("POST")
	router.HandleFunc("/api/v1/users/login", handlers.Login).Methods("POST")
	router.HandleFunc("/api/v1/users/logout", handlers.Logout).Methods("POST")
	router.Handle("/api/v1/users/mypage", common.AuthMiddleware(http.HandlerFunc(handlers.MyPageHandler))).Methods("GET")
}
