package videohub

import (
	"live/videohub/handlers"

	"github.com/gorilla/mux"
)

func RegisterRoutes(router *mux.Router) {
	videohubRouter := router.PathPrefix("/api/v1/videos").Subrouter()

	videohubRouter.HandleFunc("/list", handlers.ListVideos).Methods("GET")
	// videohubRouter.HandleFunc("/details/{id}", handlers.GetVideoDetails).Methods("GET")
}
