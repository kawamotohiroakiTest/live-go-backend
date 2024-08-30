package videoupload

import (
	"live/videoupload/handlers"

	"github.com/gorilla/mux"
)

func RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/api/v1/videoupload/upload", handlers.Upload).Methods("POST")
	// router.HandleFunc("/api/v1/videoupload/status/{id}", handlers.Status).Methods("GET")
}
