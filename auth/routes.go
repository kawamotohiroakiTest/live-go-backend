package auth

import (
	"live/auth/handlers"

	"github.com/gorilla/mux"
)

func RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/api/v1/users/register", handlers.Register).Methods("POST")
}
