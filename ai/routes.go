package ai

import (
	"live/ai/handlers"

	"github.com/gorilla/mux"
)

func RegisterRoutes(router *mux.Router) {
	recommendRouter := router.PathPrefix("/api/v1/recommend").Subrouter()

	recommendRouter.HandleFunc("/recommend", handlers.Recommend).Methods("GET")
}
